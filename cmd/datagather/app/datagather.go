package app

import (
	"context"

	"go.uber.org/zap"
	"github.com/ijsong/farseer/pkg/queue"
	"github.com/ijsong/farseer/pkg/server"
)

type DataGather struct {
	conf      *DataGatherConfig
	svr       *server.Server
	queue     *queue.EmbeddedQueue
	producers []*queue.EmbeddedQueueProducer
	consumers []*queue.EmbeddedQueueConsumer
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
		consumers: make([]*queue.EmbeddedQueueConsumer, 0),
		logger:    zap.L(),
	}

	var err error
	dg.queue, err = queue.NewEmbeddedQueue(conf.queueConfig)
	if err != nil {
		return nil, err
	}

	for i := 0; i < conf.queueConfig.NumberOfProducers; i++ {
		producer, err := queue.NewEmbeddedQueueProducer(conf.queueConfig.Address)
		if err != nil {
			return nil, err
		}
		dg.producers = append(dg.producers, producer)
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
	dg.logger.Info("starting server")
	defer dg.stop()
	return dg.svr.Start(context.Background(), services)
}

func (dg *DataGather) initService(initializer func(*queue.EmbeddedQueueProducer) (server.ServiceServer, error)) (server.ServiceServer, error) {
	var err error
	producer, err := queue.NewEmbeddedQueueProducer(dg.conf.queueConfig.Address)
	//CreateEventTopic)
	if err != nil {
		return nil, err
	}
	dg.producers = append(dg.producers, producer)

	service, err := initializer(producer)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (dg *DataGather) stop() {
	dg.logger.Info("stopping datagather")
	for _, producer := range dg.producers {
		dg.logger.Info("stopping producer")
		producer.Stop()
	}
	for _, consumer := range dg.consumers {
		dg.logger.Info("stopping consumer")
		consumer.Stop()
	}
	dg.logger.Info("stopping embedded queue")
	dg.queue.Stop()
}
