package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	prd        *kafka.Producer
	cfg        *KafkaConfig
	reportTerm chan bool
}

func NewKafkaProducer(config *KafkaConfig) (*KafkaProducer, error) {
	cmap := &kafka.ConfigMap{}
	for k, v := range config.getConfigMap() {
		cmap.SetKey(k, v)
	}
	p, err := kafka.NewProducer(cmap)
	if err != nil {
		return nil, err
	}
	prd := &KafkaProducer{
		prd:        p,
		cfg:        config,
		reportTerm: make(chan bool),
	}
	go prd.processEvents()
	return prd, nil
}

func (p *KafkaProducer) Produce(topic string, msg []byte) error {
	m := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}
	return p.prd.Produce(m, nil)
}

func (p *KafkaProducer) Close() {
	for numTry, remains := 0, 1; numTry < 5 && remains > 0; numTry++ {
		remains = p.prd.Flush(p.cfg.FlushTimeoutMs)
	}
	p.prd.Close()
	p.reportTerm <- true
}

func (p *KafkaProducer) processEvents() {
	for {
		select {
		case e := <-p.prd.Events():
			switch evt := e.(type) {
			case *kafka.Message:
				if evt.TopicPartition.Error == nil {
					continue
				}
				// handle evt => dead letter or re-enqueue
			default:
				// unknown event type
			}
		case <-p.reportTerm:
			return
		default:
		}
	}
}
