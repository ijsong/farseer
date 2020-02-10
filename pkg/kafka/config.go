package kafka

type KafkaConfig struct {
	BootstrapServers string
	FlushTimeoutMs   int
	GroupID          string
}

func NewKafkaConfig(configMap map[string]string) *KafkaConfig {
	return &KafkaConfig{}
}

func (c *KafkaConfig) getConfigMap() map[string]string {
	cmap := make(map[string]string)
	cmap["bootstrap.servers"] = c.BootstrapServers
	cmap["group.id"] = c.GroupID
	return cmap
}
