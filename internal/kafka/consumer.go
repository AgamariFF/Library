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

		var rawMessage string

		if err := json.Unmarshal(message.Value, &rawMessage); err != nil {
			logger.ErrorLog.Println("Failed to unmarshal Kafka message to string\nerr: ", err, "\nmessage.Value in string: ", string(message.Value))
			continue
		}

		var event struct {
			Data  models.Book `json:"data"`
			Event string      `json:"event"`
		}

		logger.InfoLog.Println("Raw Kafka message: ", rawMessage)
		if err := json.Unmarshal([]byte(rawMessage), &event); err != nil {
			logger.ErrorLog.Println("Failed to pars Kafka message\nerr: ", err, "\nmessage.Value in string: ", string(message.Value), "\nevent: ", event)
			continue
		}

		logger.InfoLog.Printf("New event received: %s", event.Event)

		if event.Event == "BookAdded" {
			go mailing.SendNewBookEmail(event.Data, database.DB)
		}
	}
}
