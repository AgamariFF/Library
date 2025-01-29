package kafka

import (
	"encoding/json"
	"library/internal/database"
	"library/internal/mailing"
	"library/internal/models"
	"library/logger"

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
		logger.ErrorLog.Fatal("Failed to start consumer:", err)
	}

	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		var event struct{
			event string
			data models.Book
		}
		if err := json.Unmarshal(message.Value, &event); err != nil {
			logger.ErrorLog.Println("Failed to pars Kafka message: ", err)
			continue
		}
		
		logger.InfoLog.Printf("New event received: %s", event.event)

		if event.event == "BookAdded" {
			go mailing.SendNewBookEmail(event.data, database.DB)
		}
	}
}
