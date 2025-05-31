package discord

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddFooter(t *testing.T) {
	embeds := []Embed{
		{
			Title: "Title 1",
		},
		{
			Title: "Title 2",
		},
	}

	message := Message{
		Embeds: embeds,
	}

	err := message.AddFooter("Test Footer")
	require.NoError(t, err)

	lastEmbed := message.Embeds[len(message.Embeds)-1]
	require.Equal(t, "Test Footer", lastEmbed.Footer.Text)
}

func TestAddFooterNoEmbeds(t *testing.T) {
	message := Message{}

	err := message.AddFooter("Test Footer")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrNoEmbeds)
}

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	webhook := server.URL

	message := Message{
		Embeds: []Embed{
			{
				Title: "Test Title",
			},
		},
	}

	err := message.Send(webhook)
	require.NoError(t, err)
}

func TestSendErrorStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	webhook := server.URL

	message := Message{
		Embeds: []Embed{
			{
				Title: "Test Title",
			},
		},
	}

	err := message.Send(webhook)
	require.Error(t, err)
}

func TestSendRequestError(t *testing.T) {
	message := Message{
		Embeds: []Embed{
			{
				Title: "Test Title",
			},
		},
	}

	err := message.Send("invalid-url")
	require.Error(t, err)
}

func TestSendMarshalError(t *testing.T) {
	message := Message{
		Embeds: []Embed{
			{
				Title: "Test Title",
			},
		},
	}
	// Force a marshal error by providing an invalid value
	message.Embeds[0].Color = -1

	err := message.Send("https://example.com")
	require.Error(t, err)
}
