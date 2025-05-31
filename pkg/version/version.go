package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	UpdateNotificationMessage = "[!] New update available: %s"
	GithubEndpoint            = "https://api.github.com/repos/%s/%s/releases/latest"
	errMissingGithubOwner     = errors.New("githubOwner is empty")
	errMissingGithubRepo      = errors.New("githubRepo is empty")
)

type ReleaseResponse struct {
	TagName string `json:"tag_name"`
}

type Version struct {
	GithubOwner         string
	GithubRepo          string
	LatestVersion       string
	CurrentVersion      string
	NewVersionAvailable bool
}

func (v *Version) GetUpdateNotification() string {
	if v.NewVersionAvailable && v.LatestVersion != "" {
		return fmt.Sprintf(UpdateNotificationMessage, v.LatestVersion)
	}
	return ""
}

func (v *Version) StripV() string {
	return strings.TrimPrefix(v.LatestVersion, "v")
}

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
		Timeout: 10 * time.Second,
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

func (v *Version) CheckUpdate() {
	_ = v.GetLatestVersion()
	v.NewVersionAvailable = v.CurrentVersion != v.LatestVersion
}
