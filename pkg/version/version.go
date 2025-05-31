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

// UpdateNotificationMessage is the template for update notifications.
var UpdateNotificationMessage = "[!] New update available: %s"

var githubReleaseEndpoint = "https://api.github.com/repos/%s/%s/releases/latest"

var (
	// ErrMissingGithubOwner is returned when the githubOwner is empty.
	ErrMissingGithubOwner = errors.New("githubOwner is empty")

	// ErrMissingGithubRepo is returned when the githubRepo is empty.
	ErrMissingGithubRepo = errors.New("githubRepo is empty")

	// ErrMissingCurrentVersion is returned when the currentVersion is empty.
	ErrMissingCurrentVersion = errors.New("currentVersion is empty")
)

// Versions is the interface for the version service.
type Versions interface {
	GetUpdateNotification() string
	StripV() string
	FetchLatestVersion() error
	CheckUpdate()
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
	httpClient      commonHTTPClient.Client
}

// GetUpdateNotification returns a notification string if a new version is available.
func (v *Version) GetUpdateNotification() string {
	if v.updateAvailable && v.latestVersion != "" {
		return fmt.Sprintf(UpdateNotificationMessage, v.latestVersion)
	}
	return ""
}

// StripV removes the leading 'v' from the latest version string.
func (v *Version) StripV() string {
	return strings.TrimPrefix(v.latestVersion, "v")
}

// FetchLatestVersion fetches the latest version from GitHub.
func (v *Version) FetchLatestVersion() error {
	if v.githubOwner == "" {
		return ErrMissingGithubOwner
	}
	if v.githubRepo == "" {
		return ErrMissingGithubRepo
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

	var release ReleaseResponse
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return err
	}

	v.latestVersion = release.TagName
	return nil
}

// CheckUpdate checks if a new version is available.
func (v *Version) CheckUpdate() {
	_ = v.FetchLatestVersion()
	v.updateAvailable = v.currentVersion != v.latestVersion
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

// Options is the options for the version service.
type Options struct {
	GithubOwner    string
	GithubRepo     string
	CurrentVersion string
	HTTPClient     commonHTTPClient.Client
}

func (o *Options) validate() error {
	if o.GithubOwner == "" {
		return ErrMissingGithubOwner
	}

	if o.GithubRepo == "" {
		return ErrMissingGithubRepo
	}

	if o.CurrentVersion == "" {
		return ErrMissingCurrentVersion
	}

	return nil
}

// NewVersion creates a new version service.
func NewVersion(opts Options) (Versions, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}

	if opts.HTTPClient == nil {
		opts.HTTPClient = commonHTTPClient.NewDefaultClient()
	}

	return &Version{
		githubOwner:    opts.GithubOwner,
		githubRepo:     opts.GithubRepo,
		currentVersion: opts.CurrentVersion,
		httpClient:     opts.HTTPClient,
	}, nil
}
