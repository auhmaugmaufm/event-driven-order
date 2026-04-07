package event

import (
	"fmt"
	"log"
	"time"
)

type EventType string

const (
	OrderCreated EventType = "ORDER_CREATED"
)

type Event struct {
	Type      EventType
	Payload   any
	OccuredAt time.Time
}

type Handler func(event Event)

type EventBus struct {
	subscribers map[EventType][]Handler
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventType][]Handler),
	}
}

func (eb *EventBus) Subscribe(eventType EventType, handler Handler) {
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

func (eb *EventBus) Publish(event Event) {
	handlers, ok := eb.subscribers[event.Type]
	if !ok {
		return
	}
	for _, handler := range handlers {
		go handler(event)
	}
}

func LogHandler(e Event) {
	log.Printf("[EVENT] type=%s payload=%+v occured_at=%s\n",
		e.Type,
		fmt.Sprintf("%+v", e.Payload),
		e.OccuredAt.Format(time.RFC3339),
	)
}
