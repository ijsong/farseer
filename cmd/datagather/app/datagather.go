package app

import (
	"context"
	"fmt"

	"github.com/ijsong/farseer/pkg/queue"
	"github.com/ijsong/farseer/pkg/server"
	"go.uber.org/zap"
)

type DataGather struct {
	conf      *DataGatherConfig
	svr       *server.Server
	queue     *queue.EmbeddedQueue
	producers []*queue.EmbeddedQueueProducer
	consumer  *queue.EmbeddedQueueConsumer
	logger    *zap.Logger
}

type DataGatherConfig struct {
	serverConfig *server.ServerConfig
	queueConfig  *queue.EmbeddedQueueConfig
}

type DataGatherService interface {
	Start() error
	Stop() error
}

func NewDataGather(conf *DataGatherConfig) (*DataGather, error) {
	dg := &DataGather{
		conf:      conf,
		producers: make([]*queue.EmbeddedQueueProducer, 0),
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
	return dg, nil
}

func (dg *DataGather) Start() error {
	datagatherService := NewDatagatherService(dg.producers)
	services := []server.ServiceServer{datagatherService}
	dg.queue.Start()
	dg.consumer.AddHandler(datagatherMessageHandler, dg.conf.queueConfig.NumberOfConsumers)
	dg.consumer.Connect(dg.conf.queueConfig.Address)
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
}
