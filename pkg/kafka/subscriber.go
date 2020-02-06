package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaSubscriber struct {
	csm      *kafka.Consumer
	conf     *KafkaConfig
	termchan chan interface{}
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

func (s *KafkaSubscriber) Subscribe(topic string, onSuccess func([]byte) error, onFailure func(error)) error {
	if err := s.csm.SubscribeTopics([]string{topic}, nil); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-s.termchan:
				// quit
				return
			default:
				// FIXME: set timeout parameter
				msg, err := s.csm.ReadMessage(-1)
				if err != nil {
					switch err.(kafka.Error).Code() {
					case kafka.ErrTimedOut:
						continue
					default:
						onFailure(err)
					}

				}
				onSuccess(msg.Value)

			}
		}
	}()
	return nil
}

func (s *KafkaSubscriber) Stop() error {
	s.termchan <- "quit"
	return s.csm.Close()
}
