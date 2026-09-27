package events

import (
	"encoding/json"
	"time"

	"github.com/care4u/services/booking/internal/domain"
	"github.com/google/uuid"
)

type BookingEventType string

const (
	BookingCreated    BookingEventType = "booking.created"
	BookingConfirmed  BookingEventType = "booking.confirmed"
	BookingCancelled  BookingEventType = "booking.cancelled"
	BookingCompleted  BookingEventType = "booking.completed"
	BookingNoShow     BookingEventType = "booking.no_show"
)

type BookingEvent struct {
	EventID   string          `json:"event_id"`
	EventType BookingEventType `json:"event_type"`
	Timestamp time.Time       `json:"timestamp"`
	Version   string          `json:"version"`
	Data      BookingEventData `json:"data"`
}

type BookingEventData struct {
	BookingID  string    `json:"booking_id"`
	PatientID  string    `json:"patient_id"`
	ProviderID string    `json:"provider_id"`
	ServiceType string   `json:"service_type"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes,omitempty"`
}

func NewBookingCreatedEvent(booking *domain.Booking) BookingEvent {
	return BookingEvent{
		EventID:   uuid.New().String(),
		EventType: BookingCreated,
		Timestamp: time.Now(),
		Version:   "1.0",
		Data: BookingEventData{
			BookingID:   booking.ID,
			PatientID:   booking.PatientID,
			ProviderID:  booking.ProviderID,
			ServiceType: booking.ServiceType,
			StartTime:   booking.StartTime,
			EndTime:     booking.EndTime,
			Status:      string(booking.Status),
			Notes:       booking.Notes,
		},
	}
}

func NewBookingConfirmedEvent(booking *domain.Booking) BookingEvent {
	return BookingEvent{
		EventID:   uuid.New().String(),
		EventType: BookingConfirmed,
		Timestamp: time.Now(),
		Version:   "1.0",
		Data: BookingEventData{
			BookingID:   booking.ID,
			PatientID:   booking.PatientID,
			ProviderID:  booking.ProviderID,
			ServiceType: booking.ServiceType,
			StartTime:   booking.StartTime,
			EndTime:     booking.EndTime,
			Status:      string(booking.Status),
			Notes:       booking.Notes,
		},
	}
}

func NewBookingCancelledEvent(booking *domain.Booking) BookingEvent {
	return BookingEvent{
		EventID:   uuid.New().String(),
		EventType: BookingCancelled,
		Timestamp: time.Now(),
		Version:   "1.0",
		Data: BookingEventData{
			BookingID:   booking.ID,
			PatientID:   booking.PatientID,
			ProviderID:  booking.ProviderID,
			ServiceType: booking.ServiceType,
			StartTime:   booking.StartTime,
			EndTime:     booking.EndTime,
			Status:      string(booking.Status),
			Notes:       booking.Notes,
		},
	}
}

func (e BookingEvent) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e BookingEvent) Subject() string {
	return string(e.EventType)
}