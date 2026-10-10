package events

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type EventProducer interface {
	PublishBookingCreated(event BookingEvent) error
	PublishBookingConfirmed(event BookingEvent) error
	PublishBookingCancelled(event BookingEvent) error
	Close() error
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
		Name:     "BOOKING",
		Subjects: []string{"booking.*"},
	})
	if err != nil {
		log.Printf("Stream 'BOOKING' might already exist: %v", err)
	}

	return &natsEventProducer{conn: nc, js: js}, nil
}

func (p *natsEventProducer) PublishBookingCreated(event BookingEvent) error {
	return p.publish(event)
}

func (p *natsEventProducer) PublishBookingConfirmed(event BookingEvent) error {
	return p.publish(event)
}

func (p *natsEventProducer) PublishBookingCancelled(event BookingEvent) error {
	return p.publish(event)
}

func (p *natsEventProducer) publish(event BookingEvent) error {
	payload, err := event.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	subject := event.Subject()
	_, err = p.js.Publish(subject, payload)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event: %s -> %s", subject, event.Data.BookingID)
	return nil
}

func (p *natsEventProducer) Close() error {
	p.conn.Drain()
	return nil
}
