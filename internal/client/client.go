package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/pyuldashev912/Episodes-Tracker/pkg/e"
)

const (
	getUpdateMethod           = "getUpdates"
	sendMessageMethod         = "sendMessage"
	editMessageTextMethod     = "editMessageText"
	answerCallbackQueryMethod = "answerCallbackQuery"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

// New
func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// buildParams
func buildParams(in Params) url.Values {
	out := url.Values{}

	for key, value := range in {
		out.Set(key, value)
	}

	return out
}

// Updates -> limit offset
func (c *Client) Updates(params Params) ([]Update, error) {
	const errMsg = "can't get updates"

	data, err := c.doRequest(getUpdateMethod, params)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	if !res.OK {
		return nil, e.Wrap(errMsg, err)
	}

	return res.Result, nil
}

func (c *Client) SendMessage(params Params) error {
	_, err := c.doRequest(sendMessageMethod, params)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) EditMessageText(params Params) error {
	_, err := c.doRequest(editMessageTextMethod, params)
	if err != nil {
		return e.Wrap("can't edit message", err)
	}

	return nil
}

func (c *Client) AnswerCallbackQuery(params Params) error {
	_, err := c.doRequest(answerCallbackQueryMethod, params)
	if err != nil {
		return e.Wrap("can't answer to callbak", err)
	}

	return nil
}

func (c *Client) doRequest(method string, params Params) ([]byte, error) {
	const errMsg = "can't do request"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	req.URL.RawQuery = buildParams(params).Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return body, nil
}
