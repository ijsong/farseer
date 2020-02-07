package app

import (
	"github.com/ijsong/farseer/pkg/kafka"
	"github.com/ijsong/farseer/pkg/storage"
)

type DataWriterConfig struct {
	kafkaConfig     *kafka.KafkaConfig
	cassandraConfig *storage.CassandraStorageConfig
}
