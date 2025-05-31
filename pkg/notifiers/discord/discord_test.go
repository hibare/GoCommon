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

func TestClientSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	webhook := server.URL
	client, err := NewClient(Options{
		WebhookURL: webhook,
		HTTPClient: server.Client(),
	})
	require.NoError(t, err)

	message := &Message{
		Embeds: []Embed{
			{
				Title: "Test Title",
			},
		},
	}

	err = client.Send(t.Context(), message)
	require.NoError(t, err)
}

func TestClientSendErrorStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	webhook := server.URL
	client, err := NewClient(Options{
		WebhookURL: webhook,
		HTTPClient: server.Client(),
	})
	require.NoError(t, err)

	message := &Message{
		Embeds: []Embed{
			{
				Title: "Test Title",
			},
		},
	}

	err = client.Send(t.Context(), message)
	require.Error(t, err)
}

func TestClientSendRequestError(t *testing.T) {
	client, err := NewClient(Options{
		WebhookURL: "invalid-url",
		HTTPClient: http.DefaultClient,
	})
	require.NoError(t, err)
	message := &Message{
		Embeds: []Embed{
			{
				Title: "Test Title",
			},
		},
	}
	err = client.Send(t.Context(), message)
	require.Error(t, err)
}
