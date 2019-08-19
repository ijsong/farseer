package queue

import (
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/nsqio/nsq/nsqd"
	"go.uber.org/zap"
)

type EmbeddedQueue struct {
	q *nsqd.NSQD
	c chan bool
}

type EmbeddedQueueConfig struct {
	Address           string
	DataPath          string
	MemQueueSize      int64
	SyncEvery         int64
	SyncTimeout       time.Duration
	NumberOfProducers int
	NumberOfConsumers int
}

func NewEmbeddedQueue(conf *EmbeddedQueueConfig) (*EmbeddedQueue, error) {
	opts := nsqd.NewOptions()
	opts.TCPAddress = conf.Address
	if len(conf.DataPath) > 0 {
		opts.DataPath = conf.DataPath
	}
	opts.MemQueueSize = conf.MemQueueSize
	nsqd := nsqd.New(opts)
	return &EmbeddedQueue{
		q: nsqd,
		c: make(chan bool),
	}, nil
}

func (eq *EmbeddedQueue) Start() {
	go func() {
		zap.L().Info("starting embedded queue")
		eq.q.Main()
		<-eq.c
		eq.q.Exit()
	}()
}

func (eq *EmbeddedQueue) Stop() {
	eq.c <- true
}

type EmbeddedQueueProducer struct {
	producer *nsq.Producer
}

func NewEmbeddedQueueProducer(queueAddress string) (*EmbeddedQueueProducer, error) {
	conf := nsq.NewConfig()
	producer, err := nsq.NewProducer(queueAddress, conf)
	if err != nil {
		return nil, err
	}
	return &EmbeddedQueueProducer{producer: producer}, nil
}

func (eqp *EmbeddedQueueProducer) Publish(topic string, msg []byte) error {
	return eqp.producer.Publish(topic, msg)
}

func (eqp *EmbeddedQueueProducer) Ping() error {
	return eqp.producer.Ping()
}

func (eqp *EmbeddedQueueProducer) Stop() {
	eqp.producer.Stop()
}

type EmbeddedQueueConsumer struct {
	consumer *nsq.Consumer
	topic    string
	channel  string
}

type EmbeddedQueueMessageHandler func(msg []byte) error

func NewEmbeddedQueueConsumer(topic, channel string) (*EmbeddedQueueConsumer, error) {
	conf := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, conf)
	if err != nil {
		return nil, err
	}
	return &EmbeddedQueueConsumer{consumer: consumer, topic: topic, channel: channel}, nil
}

func (eqc *EmbeddedQueueConsumer) AddHandler(handler EmbeddedQueueMessageHandler, concurrency int) {
	nsqHandler := nsq.HandlerFunc(func(m *nsq.Message) error {
		return handler(m.Body)
	})
	eqc.consumer.AddConcurrentHandlers(nsqHandler, concurrency)
}

func (eqc *EmbeddedQueueConsumer) Connect(queueAddress string) error {
	return eqc.consumer.ConnectToNSQD(queueAddress)
}

func (eqc *EmbeddedQueueConsumer) Stop() {
	eqc.consumer.Stop()
}
