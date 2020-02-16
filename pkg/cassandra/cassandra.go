package cassandra

import (
	"time"

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

func getTimestampOrNow(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now()
	}
	return t
}

func (edm *CassandraEventDataModel) Create(event *datatypes.Event) error {
	stmt, names := qb.Insert("farseer.events").Columns("user_id", "item_id", "event_type", "event_value", "version_timestamp", "event_timestamp", "properties").ToCql()
	versionTimestamp := getTimestampOrNow(event.VersionTimestamp)
	eventTimestamp := getTimestampOrNow(event.EventTimestamp)
	m := qb.M{
		"user_id":           event.UserId,
		"item_id":           event.ItemId,
		"event_type":        event.EventType,
		"event_value":       event.EventValue,
		"version_timestamp": gocql.UUIDFromTime(versionTimestamp),
		"event_timestamp":   eventTimestamp,
		"properties":        event.Properties,
	}
	q := gocqlx.Query(edm.session.Query(stmt), names).BindMap(m)
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (idm *CassandraItemDataModel) Create(item *datatypes.Item) error {
	stmt, names := qb.Insert("farseer.items").Columns("id", "item_type", "properties", "created_time", "updated_time").ToCql()
	createdTime := getTimestampOrNow(item.CreatedTime)
	updatedTime := getTimestampOrNow(item.UpdatedTime)
	m := qb.M{
		"id":           item.Id,
		"item_type":    item.ItemType,
		"properties":   item.Properties,
		"created_time": createdTime,
		"updated_time": updatedTime,
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
	stmt, names := qb.Update("farseer.items").Set("properties", "updated_time").Where(qb.Eq("id")).ToCql()
	updatedTime := getTimestampOrNow(item.UpdatedTime)
	m := qb.M{
		"id":           item.Id,
		"properties":   item.Properties,
		"updated_time": updatedTime,
	}
	q := gocqlx.Query(idm.session.Query(stmt), names).BindMap(m)
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (udm *CassandraUserDataModel) Create(user *datatypes.User) error {
	stmt, names := qb.Insert("farseer.users").Columns("id", "user_type", "properties", "created_time", "updated_time").ToCql()
	createdTime := getTimestampOrNow(user.CreatedTime)
	updatedTime := getTimestampOrNow(user.UpdatedTime)
	q := gocqlx.Query(udm.session.Query(stmt), names).BindMap(qb.M{
		"id":           user.Id,
		"user_type":    user.UserType,
		"properties":   user.Properties,
		"created_time": createdTime,
		"updated_time": updatedTime,
	})
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (udm *CassandraUserDataModel) Delete(userId string) error {
	stmt, names := qb.Delete("farseer.users").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(udm.session.Query(stmt), names).BindMap(qb.M{
		"id": userId,
	})
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}

func (udm *CassandraUserDataModel) Update(user *datatypes.User) error {
	stmt, names := qb.Update("farseer.items").Set("properties", "updated_time").Where(qb.Eq("id")).ToCql()
	updatedTime := getTimestampOrNow(user.UpdatedTime)
	m := qb.M{
		"id":           user.Id,
		"properties":   user.Properties,
		"updated_time": updatedTime,
	}
	q := gocqlx.Query(udm.session.Query(stmt), names).BindMap(m)
	if err := q.ExecRelease(); err != nil {
		zap.L().Error("could not handle query", zap.Error(err))
		return err
	}
	return nil
}
