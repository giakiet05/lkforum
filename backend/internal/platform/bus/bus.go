package bus

import (
	"sync"
)

// Event defines the interface for any event that can be published to the bus.
type Event interface {
	Topic() string
	Payload() map[string]interface{}
}

// EventListener is a channel that receives events.
type EventListener chan Event

// EventBus interface
type EventBus interface {
	Subscribe(topic string, ch EventListener)
	Publish(event Event)
}

// EventBus stores the information about subscribers, listeners and events.
type eventBus struct {
	listeners map[string][]EventListener
	lock      sync.RWMutex
}

// NewEventBus creates a new EventBus.
func NewEventBus() EventBus {
	return &eventBus{
		listeners: make(map[string][]EventListener),
	}
}

// Subscribe adds a new listener for a given topic.
func (b *eventBus) Subscribe(topic string, ch EventListener) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.listeners[topic]; !ok {
		b.listeners[topic] = make([]EventListener, 0)
	}
	b.listeners[topic] = append(b.listeners[topic], ch)
}

// Publish sends an event to all subscribed listeners of a topic.
// This is done asynchronously to prevent blocking the publisher.
func (b *eventBus) Publish(event Event) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	topic := event.Topic()
	if listeners, ok := b.listeners[topic]; ok {
		for _, listener := range listeners {
			go func(l EventListener) {
				// Use a non-blocking send to prevent a slow listener from blocking the bus.
				select {
				case l <- event:
				default:
					// Optional: Log a warning if a listener's channel is full.
				}
			}(listener)
		}
	}
}
