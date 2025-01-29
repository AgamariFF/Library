package kafka

import (
	"encoding/json"
	"errors"
	"library/logger"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.ErrorLog.Println("Failed to create Kafka producer: ", err)
		return nil, err
	}
	if producer == nil {
		logger.ErrorLog.Println("Kafka producer is nil after creation in NewKafkaProducer!")
		return nil, errors.New("producer is nil")
	}

	logger.InfoLog.Println("Kafka producer created succesfully")

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (p *KafkaProducer) SendMessage(data interface{}) error {
	messageBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(messageBytes),
	}
	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
