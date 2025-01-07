package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewKafkaConsumer(brokers []string, topic string) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (c *KafkaConsumer) ConsumeMessage() {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to start consumer:", err)
	}

	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		log.Printf("Message received: %s", string(message.Value))
	}
}
