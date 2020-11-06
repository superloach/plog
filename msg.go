package plug

import "encoding/json"

type msg struct {
	Name   string          `json:"name"`
	Call   int             `json:"call"`
	Args   json.RawMessage `json:"args,omitempty"`
	Return json.RawMessage `json:"return,omitempty"`
}
