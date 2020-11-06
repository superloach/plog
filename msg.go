package plug

import (
	"encoding/json"
)

type msg struct {
	Name   string          `json:"name"`
	Call   json.RawMessage `json:"call,omitempty"`
	Return json.RawMessage `json:"return,omitempty"`
}
