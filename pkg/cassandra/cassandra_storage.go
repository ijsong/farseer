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
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

type CassandraSession struct {
	session *gocql.Session
}

func NewCassandraStorage(conf *CassandraStorageConfig) (*CassandraStorage, error) {
	cluster := gocql.NewCluster(conf.Hosts...)
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CassandraStorage{conf: conf, cluster: cluster, session: session}, nil
}

func (s *CassandraStorage) Close() {
	s.session.Close()
}
