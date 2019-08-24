package datamodel

import (
	"github.com/gocql/gocql"
	"github.com/ijsong/farseer/pkg/datatypes"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

type EventDataModelCassandra struct {
	session *gocql.Session
}

func (e *EventDataModelCassandra) CreateEvent(event *datatypes.Event) error {
	stmt, names := qb.Insert("farseer.events").Columns("user_id", "item_id", "event_type", "event_value", "timestamp", "properties").ToCql()
	gocqlx.Query(e.session.Query(stmt), names).BindStruct(event)
	return nil
}

func (e *EventDataModelCassandra) ListEvents(userId string) ([]*datatypes.Event, error) {
	return nil, nil
}
