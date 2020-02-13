package cassandra

import (
	"github.com/gocql/gocql"
)

type CassandraStorageConfig struct {
	Hosts    []string
	Keyspace string
}

type CassandraStorage struct {
	conf    *CassandraStorageConfig
	session *gocql.Session
}

func NewCassandraStorage(conf *CassandraStorageConfig) (*CassandraStorage, error) {
	clusterConfig := gocql.NewCluster(conf.Hosts...)
	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CassandraStorage{conf: conf, session: session}, nil
}

func (s *CassandraStorage) Session() *gocql.Session {
	return s.session
}

func (s *CassandraStorage) Close() error {
	s.session.Close()
	return nil
}
