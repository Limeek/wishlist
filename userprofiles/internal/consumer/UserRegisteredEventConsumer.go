package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
	"userprofiles/internal/domain"
	"userprofiles/internal/usecase"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type UserRegisteredEvent struct {
	UserId    string `json:"userId"`
	Email     string `json:"email"`
	Phone     string `json:"phoneNumber"`
	BirthDate string `json:"birthDate"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	AvatarURL string `json:"avatarUrl"`
}

type UserRegisteredEventConsumer struct {
	kafkaReader              *kafka.Reader
	createUserProfileUseCase *usecase.CreateUserProfileUseCase
}

func NewUserRegisteredEventConsumer(createUserProfileUseCase *usecase.CreateUserProfileUseCase) *UserRegisteredEventConsumer {
	return &UserRegisteredEventConsumer{
		createUserProfileUseCase: createUserProfileUseCase,
		kafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{os.Getenv("KAFKA_BROKER")},
			Topic:          "user-registered",
			GroupID:        "user-profile-service",
			CommitInterval: 20 * time.Second,
		}),
	}
}

func (c *UserRegisteredEventConsumer) Start(ctx context.Context) {
	log.Print("User Registered Event Consumer Started")
	go func() {
		for {
			msg, err := c.kafkaReader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Println("consumer stopped")
					return
				}
				log.Printf("read error: %v", err)
				continue
			}

			c.consumeMessage(ctx, msg)
		}
	}()

}

func (c *UserRegisteredEventConsumer) consumeMessage(ctx context.Context, message kafka.Message) {
	log.Printf("Message received with key = %s", message.Key)

	if message.Value == nil {
		log.Default().Println("message value is nil")
		return
	}

	var userRegisteredEvent UserRegisteredEvent
	err := json.Unmarshal(message.Value, &userRegisteredEvent)

	if err != nil {
		log.Default().Println(err.Error())
		return
	}

	birthDate, err := time.Parse(time.DateOnly, userRegisteredEvent.BirthDate)

	if err != nil {
		log.Default().Println(err.Error())
		return
	}

	userUUID, err := uuid.Parse(userRegisteredEvent.UserId)

	if err != nil {
		log.Default().Println(err.Error())
		return
	}

	err = c.createUserProfileUseCase.CreateUserProfile(
		ctx,
		domain.UserProfile{
			Email:     userRegisteredEvent.Email,
			Phone:     userRegisteredEvent.Phone,
			LastName:  userRegisteredEvent.LastName,
			FirstName: userRegisteredEvent.FirstName,
			Avatar:    userRegisteredEvent.AvatarURL,
			BirthDate: birthDate,
			UserId:    userUUID,
		})

	if err != nil {
		log.Default().Println(err.Error())
		return
	}

	err = c.kafkaReader.CommitMessages(ctx, message)

	if err != nil {
		log.Default().Println(err.Error())
		return
	}
}

func (c *UserRegisteredEventConsumer) Close() {
	err := c.kafkaReader.Close()

	if err != nil {
		log.Default().Println("closing kafka reader error" + err.Error())
	}
}
