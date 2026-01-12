package publisher

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type UserRegisteredEventPublisher struct {
	writer *kafka.Writer
}

type UserRegisteredEvent struct {
	UserId    string `json:"userId"`
	Email     string `json:"email"`
	Phone     string `json:"phoneNumber"`
	BirthDate string `json:"birthDate"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}

func NewUserRegisteredEventPublisher() *UserRegisteredEventPublisher {
	return &UserRegisteredEventPublisher{
		&kafka.Writer{
			Addr:     kafka.TCP(os.Getenv("KAFKA_BROKER")),
			Topic:    "user-registered",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *UserRegisteredEventPublisher) PublishUserRegisteredEvent(ctx context.Context, event UserRegisteredEvent) error {
	message, err := json.Marshal(event)

	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.UserId),
		Value: message,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *UserRegisteredEventPublisher) Close() {
	err := p.writer.Close()
	if err != nil {
		log.Default().Printf("Error closing writer: %v", err)
	}
}
