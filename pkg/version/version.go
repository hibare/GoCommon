// Package version provides utilities for checking and managing application versions.
package version

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	commonHTTPClient "github.com/hibare/GoCommon/v2/pkg/http/client"
)

var (
	// UpdateNotificationMessage is the template for update notifications.
	UpdateNotificationMessage = "[!] New update available: %s"

	githubReleaseEndpoint = "https://api.github.com/repos/%s/%s/releases/latest"
)

var (
	// ErrMissingGithubOwner is returned when the githubOwner is empty.
	ErrMissingGithubOwner = errors.New("githubOwner is empty")

	// ErrMissingGithubRepo is returned when the githubRepo is empty.
	ErrMissingGithubRepo = errors.New("githubRepo is empty")

	// ErrMissingCurrentVersion is returned when the currentVersion is empty.
	ErrMissingCurrentVersion = errors.New("currentVersion is empty")

	// ErrNoTagNameInRelease is returned when the tag_name is not found in the release response.
	ErrNoTagNameInRelease = errors.New("no tag_name found in release response")
)

// VersionIface is the interface for the version service.
type VersionIface interface {
	GetUpdateNotification() string
	FetchLatestVersion() error
	CheckUpdate() error
	IsUpdateAvailable() bool
	GetLatestVersion() string
	GetCurrentVersion() string
}

// ReleaseResponse represents the response from the GitHub releases API.
type ReleaseResponse struct {
	TagName string `json:"tag_name"`
}

// Version holds version information for the application.
type Version struct {
	githubOwner     string
	githubRepo      string
	latestVersion   string
	currentVersion  string
	updateAvailable bool
	httpClient      commonHTTPClient.ClientIface
}

// IsUpdateAvailable checks if a new version is available.
func (v *Version) IsUpdateAvailable() bool {
	return v.updateAvailable
}

// GetLatestVersion returns the latest version.
func (v *Version) GetLatestVersion() string {
	return v.latestVersion
}

// GetCurrentVersion returns the current version.
func (v *Version) GetCurrentVersion() string {
	return v.currentVersion
}

func (v *Version) preChecks() error {
	if v.githubOwner == "" {
		return ErrMissingGithubOwner
	}

	if v.githubRepo == "" {
		return ErrMissingGithubRepo
	}

	return nil
}

// GetUpdateNotification returns a notification string if a new version is available.
func (v *Version) GetUpdateNotification() string {
	if v.updateAvailable && v.latestVersion != "" {
		return fmt.Sprintf(UpdateNotificationMessage, v.latestVersion)
	}
	return ""
}

// FetchLatestVersion fetches the latest version from GitHub.
func (v *Version) FetchLatestVersion() error {
	if err := v.preChecks(); err != nil {
		return err
	}

	url := fmt.Sprintf(githubReleaseEndpoint, v.githubOwner, v.githubRepo)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// Ensure we received a successful response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from GitHub releases endpoint: %d", resp.StatusCode)
	}

	var release ReleaseResponse
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return err
	}

	if release.TagName == "" {
		return ErrNoTagNameInRelease
	}

	v.latestVersion = release.TagName
	return nil
}

// CheckUpdate checks if a new version is available.
func (v *Version) CheckUpdate() error {
	if err := v.FetchLatestVersion(); err != nil {
		return err
	}

	// Compare versions with any leading "v" stripped
	current := strings.TrimPrefix(v.currentVersion, "v")
	latest := strings.TrimPrefix(v.latestVersion, "v")
	v.updateAvailable = current != latest
	return nil
}

// Options is the options for the version service.
type Options struct {
	HTTPClient commonHTTPClient.ClientIface
}

// NewVersion creates a new version service.
func NewVersion(gOwner, gRepo, cVersion string, opts Options) VersionIface {
	if opts.HTTPClient == nil {
		opts.HTTPClient = commonHTTPClient.NewDefaultClient()
	}

	return &Version{
		githubOwner:    gOwner,
		githubRepo:     gRepo,
		currentVersion: cVersion,
		httpClient:     opts.HTTPClient,
	}
}
