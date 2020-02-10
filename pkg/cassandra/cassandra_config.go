package cassandra

func (c *CassandraStorageConfig) GetHosts() []string {
	return c.Hosts
}
