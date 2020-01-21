package kafka

import (
	"fmt"
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
	p *kafka.Producer
	c *KafkaConfig
}

func NewKafkaProducer(config *KafkaConfig) (*KafkaProducer, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		return nil, err
	}
	doneChan := make(chan bool)

	go func() {
		defer close(doneChan)
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				return

			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()
	
	return &KafkaProducer{p: p, c: config}, nil
}

func (prd *KafkaProducer) Produce(topic string, msg []byte) {
	m := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}
	prd.p.ProduceChannel() <- m
}

func (prd *KafkaProducer) Close() {
	remains := prd.p.Flush(prd.c.FlushTimeoutMs)
	zap.L().Info("unflushed messages in Kafka producer", zap.Any("num", remains))
	prd.p.Close()
}

//func (kqp *KafkaQueueProducer) ProduceMessage(topic string, msg []byte) {
//	p := (*kafka.Producer)(kqp)
//	p.Produce(&kafka.Message{
//		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
//		Value:          msg,
//	}, nil)
//}
