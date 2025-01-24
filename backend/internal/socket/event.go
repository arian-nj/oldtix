package socket

import (
	"encoding/json"
)

type Event struct {
	Type EventType     `json:"type"`
	Data *EventMessage `json:"data"`
}

func EventUnmarshal(d []byte) (*Event, error) {
	event := Event{}
	err := json.Unmarshal(d, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
func (e *Event) GetJsonByte() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return []byte(""), err
	}

	return data, nil
}

func NewEvent(messagetype EventType, data EventMessage) *Event {
	e := Event{
		Type: messagetype,
		Data: &data}
	return &e
}
