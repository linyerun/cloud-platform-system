package common

import "encoding/json"

type JsonMsg struct {
	Data any    `json:"data,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func NewJsonMsgString(data any, msg string) string {
	return (&JsonMsg{Data: data, Msg: msg}).Marshal()
}

func NewJsonMsg(data []byte) *JsonMsg {
	obj := new(JsonMsg)
	obj.Unmarshal(data)
	return obj
}

func (j *JsonMsg) Marshal() string {
	bytes, _ := json.Marshal(j)
	return string(bytes)
}

func (j *JsonMsg) Unmarshal(data []byte) {
	json.Unmarshal(data, j)
}
