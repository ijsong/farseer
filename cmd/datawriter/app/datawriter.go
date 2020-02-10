package app

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/ijsong/farseer/internal/service"
	"github.com/ijsong/farseer/pkg/cassandra"
	"github.com/ijsong/farseer/pkg/datamodel"
	"github.com/ijsong/farseer/pkg/kafka"
	"github.com/ijsong/farseer/pkg/storage"
	"go.uber.org/zap"
)

type DataWriterConfig struct {
	kafkaConfig   *kafka.KafkaConfig
	storageConfig *cassandra.CassandraStorageConfig
}

type DataWriter struct {
	conf            *DataWriterConfig
	kafkaSubscriber *kafka.KafkaSubscriber
	//cassandra       *storage.CassandraStorage
	storage        storage.Storage
	eventDataModel datamodel.EventDataModel
	q              chan interface{}
}

func NewDataWriter(conf *DataWriterConfig) (*DataWriter, error) {
	ks, err := kafka.NewKafkaSubscriber(conf.kafkaConfig)
	if err != nil {
		return nil, err
	}
	cs, err := cassandra.NewCassandraStorage(conf.storageConfig)
	if err != nil {
		return nil, err
	}
	dw := &DataWriter{
		conf:            conf,
		kafkaSubscriber: ks,
		storage:         cs,
		q:               make(chan interface{}),
	}
	return dw, nil
}

func (dw *DataWriter) Start() error {
	ss, err := dw.storage.Connect()
	if err != nil {
		zap.L().Error("could not connect storage", zap.Error(err))
		return err
	}
	dw.eventDataModel, err = cassandra.NewEventDataModelCassandra(ss)
	if err != nil {
		zap.L().Error("could not instantiate data model", zap.Error(err))
		return err
	}
	dw.kafkaSubscriber.Subscribe(
		"datagather",
		func(msg []byte) error {
			// gogo cassandra
			zap.L().Info("handle message", zap.String("message", string(msg)))
			reader := bytes.NewReader(msg)
			req := &service.DatagatherRequest{}
			if err := jsonpb.Unmarshal(reader, req); err != nil {
				zap.L().Error("could not unmarshal kafka message", zap.Error(err))
				return err
			}
			fmt.Printf("req: %v\n", req)
			event := req.GetCreateEventRequest().GetEvent()
			fmt.Println(event)
			return dw.eventDataModel.CreateEvent(event)
		},
		func(err error) {
			zap.L().Error("could not handle message", zap.Error(err))
		},
	)
	ctx := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	go func() {

		defer dw.Stop()
		select {
		case <-sigs:
		case <-ctx.Done():
		}
	}()

	<-dw.q
	return nil
}

func (dw *DataWriter) Stop() error {
	err := dw.kafkaSubscriber.Stop()
	if err != nil {
		zap.L().Error("oops", zap.Error(err))
	}
	dw.q <- "term"
	return nil
}
