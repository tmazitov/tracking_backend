package ws

import "encoding/json"

type MessageType uint8

const (
	OrderUpdateStartAtFact MessageType = 1
	OrderUpdateEndAtFact   MessageType = 2
	OrderUpdateWorker      MessageType = 3
)

type OrderUpdateMessage struct {
	OrderId int64       `json:"orderId"`
	Type    MessageType `json:"type"`
	Data    interface{} `json:"data"`
}

type AuthMessage struct {
	Access string `json:"access"`
}

type Response struct {
	Status int16  `json:"status"`
	Msg    string `json:"msg"`
}

func NewResponse(status int16, msg string) *Response {
	return &Response{Status: status, Msg: msg}
}

func (e *Response) Marshal() ([]byte, error) {
	return json.Marshal(&e)
}
