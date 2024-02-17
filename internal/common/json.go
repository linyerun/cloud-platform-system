package common

import "encoding/json"

type JsonMsg struct {
	Data any    `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func NewJsonMsgString(data any, msg string) string {
	return (&JsonMsg{Data: data, Msg: msg}).Marshal()
}

func (j *JsonMsg) Marshal() string {
	bytes, _ := json.Marshal(j)
	return string(bytes)
}
