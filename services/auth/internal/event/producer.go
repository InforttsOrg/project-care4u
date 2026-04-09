package event

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/care4u/services/auth/internal/domain"
	"github.com/nats-io/nats.go"
)

type EventProducer interface {
	PublishUserVerified(user *domain.User) error
}

type natsEventProducer struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewEventProducer(url string) (EventProducer, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get jetstream context: %w", err)
	}

	// Create Stream if not exists
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "AUTH",
		Subjects: []string{"auth.*"},
	})
	if err != nil {
		log.Printf("Stream 'AUTH' might already exist: %v", err)
	}

	return &natsEventProducer{conn: nc, js: js}, nil
}

func (p *natsEventProducer) PublishUserVerified(user *domain.User) error {
	subject := "auth.user.verified"
	payload, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(subject, payload)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	log.Printf("Published event: %s -> %s", subject, user.ID)
	return nil
}
