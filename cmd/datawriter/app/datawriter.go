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
	"go.uber.org/zap"
)

type DataWriterConfig struct {
	kafkaConfig     *kafka.KafkaConfig
	cassandraConfig *cassandra.CassandraStorageConfig
}

type DataWriter struct {
	conf            *DataWriterConfig
	kafkaSubscriber *kafka.KafkaSubscriber
	cassandra       *cassandra.CassandraStorage
	userDataModel   datamodel.UserDataModel
	itemDataModel   datamodel.ItemDataModel
	eventDataModel  datamodel.EventDataModel
	q               chan interface{}
}

func NewDataWriter(conf *DataWriterConfig) (*DataWriter, error) {
	ks, err := kafka.NewKafkaSubscriber(conf.kafkaConfig)
	if err != nil {
		return nil, err
	}
	cassandra, err := cassandra.NewCassandraStorage(conf.cassandraConfig)
	if err != nil {
		return nil, err
	}
	dw := &DataWriter{
		conf:            conf,
		kafkaSubscriber: ks,
		cassandra:       cassandra,
		q:               make(chan interface{}),
	}
	return dw, nil
}

func (dw *DataWriter) Start() error {
	var err error
	dw.eventDataModel, err = cassandra.NewCassandraEventDataModel(dw.cassandra)
	if err != nil {
		zap.L().Error("could not instantiate data model", zap.Error(err))
		return err
	}
	dw.itemDataModel, err = cassandra.NewCassandraItemDataModel(dw.cassandra)
	if err != nil {
		zap.L().Error("could not instantiate data model", zap.Error(err))
		return err
	}
	dw.userDataModel, err = cassandra.NewCassandraUserDataModel(dw.cassandra)
	if err != nil {
		zap.L().Error("could not instantiate data model", zap.Error(err))
		return err
	}

	dw.kafkaSubscriber.Subscribe(
		"datagather",
		func(msg []byte) error {
			req, err := dw.parseDatagatherRequest(msg)
			if err != nil {
				return err
			}
			return dw.handleDatagatherRequest(req)
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

	zap.L().Info("starting...")

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

func (dw *DataWriter) parseDatagatherRequest(msg []byte) (*service.DatagatherRequest, error) {
	// gogo cassandra
	zap.L().Info("handle message", zap.String("message", string(msg)))
	reader := bytes.NewReader(msg)
	req := &service.DatagatherRequest{}
	if err := jsonpb.Unmarshal(reader, req); err != nil {
		zap.L().Error("could not unmarshal kafka message", zap.Error(err))
		return nil, err
	}
	return req, nil
}

func (dw *DataWriter) handleDatagatherRequest(req *service.DatagatherRequest) error {
	switch msg := req.GetValue().(type) {
	case *service.CreateEventRequest:
		return dw.eventDataModel.Create(msg.GetEvent())
	case *service.CreateItemRequest:
		return dw.itemDataModel.Create(msg.GetItem())
	case *service.DeleteItemRequest:
		return dw.itemDataModel.Delete(msg.GetItemId())
	case *service.UpdateItemRequest:
		return dw.itemDataModel.Update(msg.GetItem())
	case *service.CreateUserRequest:
		return dw.userDataModel.Create(msg.GetUser())
	case *service.DeleteUserRequest:
		return dw.userDataModel.Delete(msg.GetUserId())
	case *service.UpdateUserRequest:
		return dw.userDataModel.Update(msg.GetUser())
	default:
		zap.L().Error("could not handle message")
		return fmt.Errorf("unknown request message")
	}
}
