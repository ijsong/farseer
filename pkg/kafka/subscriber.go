package kafka

import "github.com/confluentinc/confluent-kafka-go/kafka"

type KafkaSubscriber struct {
	csm  *kafka.Consumer
	conf *KafkaConfig
}

func NewKafkaSubscriber(conf *KafkaConfig) (*KafkaSubscriber, error) {
	cmap := &kafka.ConfigMap{}
	for k, v := range conf.getConfigMap() {
		cmap.SetKey(k, v)
	}
	csm, err := kafka.NewConsumer(cmap)
	if err != nil {
		return nil, err
	}
	sub := &KafkaSubscriber{
		csm:  csm,
		conf: conf,
	}
	return sub, nil
}

func (s *KafkaSubscriber) Subscribe(topics []string) error {
	return nil
}
