package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/ijsong/farseer/pkg/datatypes"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"go.uber.org/zap"
)

type CassandraEventDataModel struct {
	session *gocql.Session
}

type CassandraItemDataModel struct {
	session *gocql.Session
}

type CassandraUserDataModel struct {
	session *gocql.Session
}

func NewCassandraEventDataModel(s *CassandraStorage) (*CassandraEventDataModel, error) {
	return &CassandraEventDataModel{
		session: s.Session(),
	}, nil
}

func NewCassandraItemDataModel(s *CassandraStorage) (*CassandraItemDataModel, error) {
	return &CassandraItemDataModel{
		session: s.Session(),
	}, nil
}

func NewCassandraUserDataModel(s *CassandraStorage) (*CassandraUserDataModel, error) {
	return &CassandraUserDataModel{
		session: s.Session(),
	}, nil
}

func (edm *CassandraEventDataModel) Create(event *datatypes.Event) error {
	stmt, names := qb.Insert("farseer.events").Columns("user_id", "item_id", "event_type", "event_value", "timestamp", "properties").ToCql()
	m := qb.M{
		"user_id":     event.UserId,
		"item_id":     event.ItemId,
		"event_type":  event.EventType,
		"event_value": event.EventValue,
		"timestamp":   gocql.UUIDFromTime(event.Timestamp),
		"properties":  event.Properties,
	}
	q := gocqlx.Query(edm.session.Query(stmt), names).BindMap(m)
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (idm *CassandraItemDataModel) Create(item *datatypes.Item) error {
	stmt, names := qb.Insert("farseer.items").Columns("id", "properties", "create_time", "update_time").ToCql()
	m := qb.M{
		"id":          item.Id,
		"properties":  item.Properties,
		"create_time": gocql.UUIDFromTime(item.CreateTime),
		"update_time": gocql.UUIDFromTime(item.UpdateTime),
	}
	q := gocqlx.Query(idm.session.Query(stmt), names).BindMap(m)
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (idm *CassandraItemDataModel) Delete(itemId string) error {
	stmt, names := qb.Delete("farseer.items").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(idm.session.Query(stmt), names).BindMap(qb.M{
		"id": itemId,
	})
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (idm *CassandraItemDataModel) Update(item *datatypes.Item) error {
	stmt, names := qb.Update("farseer.items").Set("properties", "update_time").Where(qb.Eq("id")).ToCql()
	m := qb.M{
		"id":          item.Id,
		"properties":  item.Properties,
		"update_time": gocql.UUIDFromTime(item.UpdateTime),
	}
	q := gocqlx.Query(idm.session.Query(stmt), names).BindMap(m)
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (udm *CassandraUserDataModel) Create(user *datatypes.User) error {
	stmt, names := qb.Insert("farseer.users").Columns("id", "properties", "create_time", "update_time").ToCql()
	q := gocqlx.Query(udm.session.Query(stmt), names).BindMap(qb.M{
		"id":          user.Id,
		"properties":  user.Properties,
		"create_time": gocql.UUIDFromTime(user.CreateTime),
		"update_time": gocql.UUIDFromTime(user.UpdateTime),
	})
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (udm *CassandraUserDataModel) Delete(userId string) error {
	return nil
}

func (udm *CassandraUserDataModel) Update(user *datatypes.User) error {
	return nil
}
