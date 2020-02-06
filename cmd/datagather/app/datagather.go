package app

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/ijsong/farseer/internal/service"
	"github.com/ijsong/farseer/pkg/kafka"

	"github.com/ijsong/farseer/pkg/queue"
	"github.com/ijsong/farseer/pkg/server"
	"github.com/ijsong/farseer/pkg/storage"
	"go.uber.org/zap"
)

type DataGather struct {
	conf          *DataGatherConfig
	svr           *server.Server
	queue         *queue.EmbeddedQueue
	producers     []*queue.EmbeddedQueueProducer
	consumer      *queue.EmbeddedQueueConsumer
	kafkaProducer *kafka.KafkaProducer
	storage       storage.Storage
	logger        *zap.Logger
}

type DataGatherConfig struct {
	serverConfig *server.ServerConfig
	queueConfig  *queue.EmbeddedQueueConfig
	kafkaConfig  *kafka.KafkaConfig
}

type DataGatherService interface {
	Start() error
	Stop() error
}

func NewDataGather(conf *DataGatherConfig) (*DataGather, error) {
	dg := &DataGather{
		conf:      conf,
		producers: nil,
		consumer:  nil,
		logger:    zap.L(),
	}

	var err error
	dg.queue, err = queue.NewEmbeddedQueue(conf.queueConfig)
	if err != nil {
		return nil, err
	}

	for i := 0; i < conf.queueConfig.NumberOfProducers; i++ {
		producer, err := queue.NewEmbeddedQueueProducer(conf.queueConfig.Address, datagatherTopic)
		if err != nil {
			return nil, err
		}
		dg.producers = append(dg.producers, producer)
	}

	channel := fmt.Sprintf("%s_channel", datagatherTopic)
	dg.consumer, err = queue.NewEmbeddedQueueConsumer(datagatherTopic, channel)
	if err != nil {
		return nil, err
	}

	dg.svr, err = server.NewServer(conf.serverConfig)
	if err != nil {
		return nil, err
	}

	dg.storage = storage.NewCassandraStorage(strings.Split(conf.cassandraHosts, ","))

	dg.kafkaProducer, err = kafka.NewKafkaProducer(conf.kafkaConfig)
	if err != nil {
		return nil, err
	}
	return dg, nil
}

func (dg *DataGather) Start() error {
	datagatherService := NewDatagatherService(dg.producers)
	services := []server.ServiceServer{datagatherService}
	dg.queue.Start()
	dg.consumer.AddHandler(func(msg []byte) error {
		req := &service.DatagatherRequest{}
		if err := proto.Unmarshal(msg, req); err != nil {
			return err
		}

		marshaler := &jsonpb.Marshaler{}
		var buf bytes.Buffer
		if err := marshaler.Marshal(&buf, req); err != nil {
			return err
		}
		bytes := buf.Bytes()

		zap.L().Info("message handler", zap.Any("req", buf.String()))
		return dg.kafkaProducer.Produce("datagather", bytes)
	}, dg.conf.queueConfig.NumberOfConsumers)
	dg.consumer.Connect(dg.conf.queueConfig.Address)
	dg.storage.Connect()
	dg.logger.Info("starting server")
	defer dg.stop()
	return dg.svr.Start(context.Background(), services)
}

func (dg *DataGather) stop() {
	dg.logger.Info("stopping datagather")
	for _, producer := range dg.producers {
		dg.logger.Info("stopping producer")
		producer.Stop()
	}
	dg.logger.Info("stopping consumer")
	dg.consumer.Stop()

	dg.logger.Info("stopping embedded queue")
	dg.queue.Stop()

	dg.logger.Info("stopping kafka producer")
	dg.kafkaProducer.Close()
}
