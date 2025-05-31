package discord

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	lastEmbed := message.Embeds[len(message.Embeds)-1]
	assert.Equal(t, "Test Footer", lastEmbed.Footer.Text)
}

func TestAddFooterNoEmbeds(t *testing.T) {
	message := Message{}

	err := message.AddFooter("Test Footer")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNoEmbeds)
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
	assert.NoError(t, err)
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
	assert.Error(t, err)
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
	assert.Error(t, err)
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
	assert.Error(t, err)
}
