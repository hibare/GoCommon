package version

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const expectedJSON = `{
		"url": "https://api.github.com/repos/hibare/Sample/releases/61364471",
		"assets_url": "https://api.github.com/repos/hibare/Sample/releases/61364471/assets",
		"upload_url": "https://uploads.github.com/repos/hibare/Sample/releases/61364471/assets{?name,label}",
		"html_url": "https://github.com/hibare/Sample/releases/tag/v0.7.1",
		"id": 61364471,
		"author": {
			"login": "hibare",
			"id": 20609766,
			"node_id": "MDQ6VXNlcjIwNjA5NzY2",
			"avatar_url": "https://avatars.githubusercontent.com/u/20609766?v=4",
			"gravatar_id": "",
			"url": "https://api.github.com/users/hibare",
			"html_url": "https://github.com/hibare",
			"followers_url": "https://api.github.com/users/hibare/followers",
			"following_url": "https://api.github.com/users/hibare/following{/other_user}",
			"gists_url": "https://api.github.com/users/hibare/gists{/gist_id}",
			"starred_url": "https://api.github.com/users/hibare/starred{/owner}{/repo}",
			"subscriptions_url": "https://api.github.com/users/hibare/subscriptions",
			"organizations_url": "https://api.github.com/users/hibare/orgs",
			"repos_url": "https://api.github.com/users/hibare/repos",
			"events_url": "https://api.github.com/users/hibare/events{/privacy}",
			"received_events_url": "https://api.github.com/users/hibare/received_events",
			"type": "User",
			"site_admin": false
		},
		"node_id": "RE_kwDODVj4p84DqFj3",
		"tag_name": "v0.7.1",
		"target_commitish": "master",
		"name": "v0.7.1",
		"draft": false,
		"prerelease": false,
		"created_at": "2022-03-09T12:00:07Z",
		"published_at": "2022-03-09T12:00:38Z",
		"assets": [

		],
		"tarball_url": "https://api.github.com/repos/hibare/Sample/tarball/v0.7.1",
		"zipball_url": "https://api.github.com/repos/hibare/Sample/zipball/v0.7.1",
		"body": "## Something juicy changes",
		"mentions_count": 2
	}`

func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.github.v3+json")
		_, _ = w.Write([]byte(expectedJSON))
	}))
	defer server.Close()

	GithubEndpoint = server.URL + "#%s#%s"

	code := m.Run()
	os.Exit(code)
}

func TestCheckLatestVersion(t *testing.T) {
	version := Version{
		GithubOwner: "hibare",
		GithubRepo:  "Sample",
	}
	err := version.GetLatestVersion()
	assert.NoError(t, err)
	assert.Equal(t, "v0.7.1", version.LatestVersion)
}

func TestLatestVersionMissingGitOwner(t *testing.T) {
	version := Version{}
	err := version.GetLatestVersion()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errMissingGithubOwner)
}

func TestLatestVersionMissingGitRepo(t *testing.T) {
	version := Version{
		GithubOwner: "hibare",
	}
	err := version.GetLatestVersion()
	assert.Error(t, err)
	assert.ErrorIs(t, err, errMissingGithubRepo)
}

func TestIsNewVersionAvailableTrue(t *testing.T) {
	version := Version{
		GithubOwner:    "hibare",
		GithubRepo:     "Sample",
		CurrentVersion: "0.0.0",
	}
	version.CheckUpdate()
	assert.True(t, version.NewVersionAvailable)
	assert.Equal(t, "v0.7.1", version.LatestVersion)
}

func TestIsNewVersionAvailableFalse(t *testing.T) {
	version := Version{
		GithubOwner:    "hibare",
		GithubRepo:     "Sample",
		CurrentVersion: "v0.7.1",
	}
	version.CheckUpdate()
	assert.False(t, version.NewVersionAvailable)
	assert.Equal(t, "v0.7.1", version.LatestVersion)
}

func TestIsNewVersionAvailableFailure(t *testing.T) {
	version := Version{}
	version.CheckUpdate()
	assert.False(t, version.NewVersionAvailable)
}

func TestGetUpdateNotification(t *testing.T) {
	version := Version{
		LatestVersion: "0.7.1",
	}
	version.CheckUpdate()
	assert.Equal(t, "[!] New update available: 0.7.1", version.GetUpdateNotification())
}

func TestGetUpdateNotificationNoUpdate(t *testing.T) {
	version := Version{
		LatestVersion: "0.7.1",
	}
	assert.Empty(t, version.GetUpdateNotification())
}
