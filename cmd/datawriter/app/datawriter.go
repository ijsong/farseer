package app

import (
	"github.com/ijsong/farseer/pkg/kafka"
	"go.uber.org/zap"
)

type DataWriter struct {
	conf            *DataWriterConfig
	kafkaSubscriber *kafka.KafkaSubscriber
}

func NewDataWriter(conf *DataWriterConfig) (*DataWriter, error) {
	ks, err := kafka.NewKafkaSubscriber(conf.kafkaConfig)
	if err != nil {
		return nil, err
	}
	dw := &DataWriter{
		conf:            conf,
		kafkaSubscriber: ks,
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
