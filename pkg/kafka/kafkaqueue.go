package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type KafkaQueue struct {
}

type KafkaConfig struct {
	BootstrapServers string
	FlushTimeoutMs   int
}

type KafkaProducer struct {
	p          *kafka.Producer
	c          *KafkaConfig
	reportTerm chan bool
}

func NewKafkaProducer(config *KafkaConfig) (*KafkaProducer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, err
	}
	prd := &KafkaProducer{p: p, c: config, reportTerm: make(chan bool)}
	go prd.processEvents()
	return prd, nil
}

func (prd *KafkaProducer) Produce(topic string, msg []byte) error {
	m := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}
	return prd.p.Produce(m, nil)
	//prd.p.ProduceChannel() <- m
}

func (prd *KafkaProducer) Close() {
	remains := prd.p.Flush(prd.c.FlushTimeoutMs)
	zap.L().Info("unflushed messages in Kafka producer", zap.Any("num", remains))
	prd.p.Close()
	prd.reportTerm <- true
}

func (prd *KafkaProducer) processEvents() {
	for {
		select {
		case e := <-prd.p.Events():
			switch evt := e.(type) {
			case *kafka.Message:
				if evt.TopicPartition.Error == nil {
					continue
				}
				// handle evt => dead letter or re-enqueue
			default:
				// unknown event type
			}
		case <-prd.reportTerm:
			return
		default:
		}
	}
}
