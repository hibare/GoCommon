package version

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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

var testHTTPClient *http.Client

func TestMain(m *testing.M) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.github.v3+json")
		_, _ = w.Write([]byte(expectedJSON))
	}))
	githubReleaseEndpoint = server.URL + "#%s#%s"
	testHTTPClient = server.Client()

	code := m.Run()

	server.Close()
	os.Exit(code)
}

func TestCheckLatestVersion(t *testing.T) {
	version := Version{
		githubOwner: "hibare",
		githubRepo:  "Sample",
		httpClient:  testHTTPClient,
	}
	err := version.FetchLatestVersion()
	require.NoError(t, err)
	require.Equal(t, "v0.7.1", version.latestVersion)
}

func TestLatestVersionMissingGitOwner(t *testing.T) {
	version := Version{httpClient: testHTTPClient}
	err := version.FetchLatestVersion()
	require.Error(t, err)
	require.ErrorIs(t, err, ErrMissingGithubOwner)
}

func TestLatestVersionMissingGitRepo(t *testing.T) {
	version := Version{
		githubOwner: "hibare",
		httpClient:  testHTTPClient,
	}
	err := version.FetchLatestVersion()
	require.Error(t, err)
	require.ErrorIs(t, err, ErrMissingGithubRepo)
}

func TestIsNewVersionAvailableTrue(t *testing.T) {
	version := Version{
		githubOwner:    "hibare",
		githubRepo:     "Sample",
		currentVersion: "0.0.0",
		httpClient:     testHTTPClient,
	}
	version.CheckUpdate()
	require.True(t, version.updateAvailable)
	require.Equal(t, "v0.7.1", version.latestVersion)
}

func TestIsNewVersionAvailableFalse(t *testing.T) {
	version := Version{
		githubOwner:    "hibare",
		githubRepo:     "Sample",
		currentVersion: "v0.7.1",
		httpClient:     testHTTPClient,
	}
	version.CheckUpdate()
	require.False(t, version.updateAvailable)
	require.Equal(t, "v0.7.1", version.latestVersion)
}

func TestIsNewVersionAvailableFailure(t *testing.T) {
	version := Version{}
	version.CheckUpdate()
	require.False(t, version.updateAvailable)
}

func TestGetUpdateNotification(t *testing.T) {
	version := Version{
		latestVersion: "0.7.1",
	}
	version.CheckUpdate()
	require.Equal(t, "[!] New update available: 0.7.1", version.GetUpdateNotification())
}

func TestGetUpdateNotificationNoUpdate(t *testing.T) {
	version := Version{
		latestVersion: "0.7.1",
	}
	require.Empty(t, version.GetUpdateNotification())
}

func TestStripV(t *testing.T) {
	v := Version{latestVersion: "v1.2.3"}
	require.Equal(t, "1.2.3", strings.TrimPrefix(v.latestVersion, "v"))
	v2 := Version{latestVersion: "1.2.3"}
	require.Equal(t, "1.2.3", strings.TrimPrefix(v2.latestVersion, "v"))
}

func TestGetCurrentVersion(t *testing.T) {
	v := Version{currentVersion: "v0.1.0"}
	require.Equal(t, "v0.1.0", v.GetCurrentVersion())
}

func TestIsUpdateAvailableDirect(t *testing.T) {
	v := Version{updateAvailable: true}
	require.True(t, v.IsUpdateAvailable())
	v2 := Version{updateAvailable: false}
	require.False(t, v2.IsUpdateAvailable())
}

func TestFetchLatestVersion_ErrorCases(t *testing.T) {
	v := Version{}
	err := v.FetchLatestVersion()
	require.ErrorIs(t, err, ErrMissingGithubOwner)

	v = Version{githubOwner: "owner"}
	err = v.FetchLatestVersion()
	require.ErrorIs(t, err, ErrMissingGithubRepo)
}

func TestNewVersion_Success(t *testing.T) {
	opts := Options{}
	vs := NewVersion("owner", "repo", "v1.0.0", opts)
	require.NotNil(t, vs)
}

func TestNewVersion_DefaultHTTPClient(t *testing.T) {
	opts := Options{
		HTTPClient: nil,
	}
	vs := NewVersion("owner", "repo", "v1.0.0", opts)
	ver, ok := vs.(*Version)
	require.True(t, ok)
	require.NotNil(t, ver.httpClient)
}

func TestVersionServiceInterface(t *testing.T) {
	opts := Options{}
	_ = NewVersion("owner", "repo", "v1.0.0", opts)
}
