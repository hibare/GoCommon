package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	CurrentVersion            = "0.0.0"
	LatestVersion             = CurrentVersion
	UpdateNotificationMessage = "New update available: %s"
	GithubEndpoint            = "https://api.github.com/repos/%s/%s/releases/latest"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func CheckLatestRelease(githubOwner, githubRepo string) {
	url := fmt.Sprintf(GithubEndpoint, githubOwner, githubRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		fmt.Print(err)
		return
	}

	LatestVersion = strings.TrimPrefix(release.TagName, "v")
}

func IsNewVersionAvailable() bool {
	return CurrentVersion != LatestVersion
}

func GetUpdateNotification() string {
	return fmt.Sprintf(UpdateNotificationMessage, LatestVersion)
}
