package storage

import "github.com/gocql/gocql"

type CassandraStorage struct {
	conf    *CassandraStorageConfig
	cluster *gocql.ClusterConfig
}

type CassandraSession struct {
	session *gocql.Session
}

func NewCassandraStorage(conf *CassandraStorageConfig) (*CassandraStorage, error) {
	cluster := gocql.NewCluster(conf.Hosts...)
	return &CassandraStorage{cluster: cluster}, nil
}

func (c *CassandraStorage) Connect() (StorageSession, error) {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CassandraSession{session: session}, nil
}

func (s *CassandraSession) GetUnderlying() *gocql.Session {
	return s.session
}

func (s *CassandraSession) Close() {
	s.Close()
}
