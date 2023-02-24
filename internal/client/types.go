package client

import (
	"encoding/json"
	"strconv"

	"github.com/pyuldashev912/tracker/pkg/e"
)

// Params
type Params map[string]string

// AddParam
func (p Params) AddParam(key string, value interface{}) error {
	switch v := value.(type) {
	case int:
		p[key] = strconv.Itoa(v)
	case string:
		p[key] = v
	default:
		result, err := json.Marshal(value)
		if err != nil {
			e.Wrap("can't marshal", err)
		}

		p[key] = string(result)
	}

	return nil
}

// User
type User struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
}

// Chat
type Chat struct {
	ID int `json:"id"`
}

// IncomingMessage
type IncomingMessage struct {
	Text string `json:"text"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
}

// Update
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

// UpdateResponse
type UpdateResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}
