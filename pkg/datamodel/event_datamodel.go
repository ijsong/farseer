package datamodel

import "github.com/ijsong/farseer/pkg/datatypes"

type EventDataModel interface {
	CreateEvent(event *datatypes.Event) error
	ListEvents(userId string) ([]*datatypes.Event, error)
}
