package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/ijsong/farseer/pkg/datatypes"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"go.uber.org/zap"
)

type EventDataModelCassandra struct {
	storage *CassandraStorage
}

func NewEventDataModelCassandra(s *CassandraStorage) (*EventDataModelCassandra, error) {
	return &EventDataModelCassandra{
		storage: s,
	}, nil
}

func (e *EventDataModelCassandra) CreateEvent(event *datatypes.Event) error {
	stmt, names := qb.Insert("farseer.events").Columns("user_id", "item_id", "event_type", "event_value", "timestamp", "properties").ToCql()
	m := toMap(event)
	q := gocqlx.Query(e.storage.session.Query(stmt), names).BindMap(m)
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (e *EventDataModelCassandra) ListEvents(userId string) ([]*datatypes.Event, error) {
	return nil, nil
}

func toMap(event *datatypes.Event) map[string]interface{} {
	return qb.M{
		"user_id":     event.UserId,
		"item_id":     event.ItemId,
		"event_type":  event.EventType,
		"event_value": event.EventValue,
		"timestamp":   gocql.UUIDFromTime(event.Timestamp),
		"properties":  event.Properties,
	}

}
