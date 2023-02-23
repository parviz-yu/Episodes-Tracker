package client

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestClient_Update(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		jsonAnswer := `{"ok":true,"result":[{"update_id":1,"message":{"text":"mock text","from":{"username":"mock username","first_name":"mock first name"},"chat":{"id":1}}}]}`
		assert.Equal(t, req.URL.String(), "https://api.telegram.org/bottoken/getUpdate?limit=100&offset=0")
		return &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(jsonAnswer)),
		}
	})

	c := Client{
		host:     "api.telegram.org",
		basePath: "bottoken",
		client:   *client,
	}

	p := Params{}
	p.AddParam("offset", 0)
	p.AddParam("limit", 100)

	data, err := c.Update(p)
	assert.NoError(t, err)
	assert.NotNil(t, data)
}

func TestClient_SendMessage(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), "https://api.telegram.org/bottoken/sendMessage?limit=100&offset=0")
		return &http.Response{}
	})

	c := Client{
		host:     "api.telegram.org",
		basePath: "bottoken",
		client:   *client,
	}

	p := Params{}
	p.AddParam("offset", 0)
	p.AddParam("limit", 100)

	err := c.SendMessage(p)
	assert.NoError(t, err)
}
