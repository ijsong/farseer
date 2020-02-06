package storage

import "github.com/gocql/gocql"

type CassandraStorage struct {
	conf    *CassandraStorageConfig
	cluster *gocql.ClusterConfig
}

type CassandraSession struct {
	session *gocql.Session
}

func NewCassandraStorage(conf *CassandraStorageConfig) Storage {
	cluster := gocql.NewCluster(conf.Hosts...)
	return &CassandraStorage{cluster: cluster}
}

func (c *CassandraStorage) Connect() (StorageSession, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CassandraSession{session: session}, nil
}

func (s *CassandraSession) Close() {
	s.Close()
}
