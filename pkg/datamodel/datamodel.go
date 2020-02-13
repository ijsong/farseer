package datamodel

import (
	"github.com/ijsong/farseer/pkg/datatypes"
)

type EventDataModel interface {
	Create(event *datatypes.Event) error
	// GetEvents(userId string) ([]*datatypes.Event, error)
}

type ItemDataModel interface {
	Create(item *datatypes.Item) error
	Delete(itemId string) error
	Update(item *datatypes.Item) error
	// GetItem(itemId string) (*datatypes.Item, error)
}

type UserDataModel interface {
	Create(user *datatypes.User) error
	Delete(userId string) error
	Update(user *datatypes.User) error
	// GetUser(userId string) (*datatypes.User, error)
}
