package storage

import "github.com/gocql/gocql"

type CassandraStorage struct {
	cluster *gocql.ClusterConfig
}

type CassandraSession struct {
	session *gocql.Session
}

func NewCassandraStorage(hosts []string) Storage {
	cluster := gocql.NewCluster(hosts...)
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
