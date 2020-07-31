package pubsub

import (
	"encoding/json"
	"io"
)

type Message struct {
	Subscription string        `json:subscription`
	Message      Base64Message `json:"message"`
	MessageID    string        `json:"messageID"`
}
type Base64Message struct {
	Data []byte `json:"data"`
}

func DecodePubSubMessage(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, nil
	}

	var event Message

	err := json.NewDecoder(r).Decode(&event)
	if err != nil {
		return nil, err
	}

	if event.Message.Data == nil {
		return nil, nil
	}

	return []byte(string(event.Message.Data)), nil

}
