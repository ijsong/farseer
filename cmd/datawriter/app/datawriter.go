package app

import (
	"github.com/ijsong/farseer/internal/datamodel"
	"github.com/ijsong/farseer/pkg/kafka"
	"github.com/ijsong/farseer/pkg/storage"
	"go.uber.org/zap"
)

type DataWriter struct {
	conf            *DataWriterConfig
	kafkaSubscriber *kafka.KafkaSubscriber
	cassandra       *storage.CassandraStorage
	eventDataModel  datamodel.EventDataModel
}

func NewDataWriter(conf *DataWriterConfig) (*DataWriter, error) {
	ks, err := kafka.NewKafkaSubscriber(conf.kafkaConfig)
	if err != nil {
		return nil, err
	}
	cs, err := storage.NewCassandraStorage(conf.cassandraConfig)
	if err != nil {
		return nil, err
	}
	dw := &DataWriter{
		conf:            conf,
		kafkaSubscriber: ks,
		cassandra:       cs,
	}
	return dw, nil
}

func (dw *DataWriter) Start() error {
	dw.kafkaSubscriber.Subscribe(
		"datagather",
		func(msg []byte) error {
			// gogo cassandra
			zap.L().Info("handle message", zap.String("message", string(msg)))
			return nil
		},
		func(err error) {
			zap.L().Error("could not handle message", zap.Error(err))
		},
	)
	return nil
}

func (dw *DataWriter) Stop() error {
	return nil
}
