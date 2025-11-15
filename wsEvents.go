package main

import "errors"

type WsEvent struct {
	Event           string                 `json:"event"`
	EventMsgDetails map[string]interface{} `json:"-"`
}

type ClientSendingEvents int

const (
	mmQueue ClientSendingEvents = iota
	mmUnqueue
	mWordSubmission
	mWordChange
)

func GetClientEventFromStr(event string) (ClientSendingEvents, error) {
	switch event {
	case "mm.queue":
		return mmQueue, nil
	case "mm.unqueue":
		return mmUnqueue, nil
	case "m.word-submission":
		return mWordSubmission, nil
	case "m.word-change":
		return mWordChange, nil
	}

	return -1, errors.New("client event does not exist")

}
