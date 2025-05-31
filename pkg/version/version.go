// Package version provides utilities for checking and managing application versions.
package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	commonHTTP "github.com/hibare/GoCommon/v2/pkg/http"
)

// UpdateNotificationMessage is the template for update notifications.
var UpdateNotificationMessage = "[!] New update available: %s"

// GithubEndpoint is the GitHub API endpoint for latest releases.
var GithubEndpoint = "https://api.github.com/repos/%s/%s/releases/latest"

var (
	errMissingGithubOwner = errors.New("githubOwner is empty")
	errMissingGithubRepo  = errors.New("githubRepo is empty")
)

// ReleaseResponse represents the response from the GitHub releases API.
type ReleaseResponse struct {
	TagName string `json:"tag_name"`
}

// Version holds version information for the application.
type Version struct {
	GithubOwner         string
	GithubRepo          string
	LatestVersion       string
	CurrentVersion      string
	NewVersionAvailable bool
}

// GetUpdateNotification returns a notification string if a new version is available.
func (v *Version) GetUpdateNotification() string {
	if v.NewVersionAvailable && v.LatestVersion != "" {
		return fmt.Sprintf(UpdateNotificationMessage, v.LatestVersion)
	}
	return ""
}

// StripV removes the leading 'v' from the latest version string.
func (v *Version) StripV() string {
	return strings.TrimPrefix(v.LatestVersion, "v")
}

// GetLatestVersion fetches the latest version from GitHub.
func (v *Version) GetLatestVersion() error {
	if v.GithubOwner == "" {
		return errMissingGithubOwner
	}

	if v.GithubRepo == "" {
		return errMissingGithubRepo
	}

	url := fmt.Sprintf(GithubEndpoint, v.GithubOwner, v.GithubRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{
		Timeout: commonHTTP.DefaultHTTPClientTimeout,
	}
	resp, err := client.Do(req)
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

	v.LatestVersion = release.TagName
	return nil
}

// CheckUpdate checks if a new version is available.
func (v *Version) CheckUpdate() {
	_ = v.GetLatestVersion()
	v.NewVersionAvailable = v.CurrentVersion != v.LatestVersion
}
