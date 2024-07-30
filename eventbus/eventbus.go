package eventbus

import (
	"sync"
)

type Event struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type EventBus struct {
	subscribers map[string]map[string]chan Event
	mu          sync.RWMutex
}

var (
	instance *EventBus
	once     sync.Once
	mutex    sync.Mutex
)

func NewEventBus() *EventBus {
	if instance == nil {
		mutex.Lock()
		defer mutex.Unlock()
		if instance == nil {
			once.Do(func() {
				instance = &EventBus{
					subscribers: make(map[string]map[string]chan Event),
				}
			})
		}
	}
	return instance
}

func (eb *EventBus) Subscribe(id string, eventName string, ch chan Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if _, ok := eb.subscribers[eventName]; !ok {
		eb.subscribers[eventName] = make(map[string]chan Event)
	}
	if _, ok := eb.subscribers[eventName][id]; !ok {
		eb.subscribers[eventName][id] = ch
	} else {
		eb.subscribers[eventName][id] = ch
	}
}

func (eb *EventBus) Unsubscribe(id string, eventName string) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	if _, ok := eb.subscribers[eventName][id]; !ok {
		return
	}
	delete(eb.subscribers[eventName], id)
}

func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if chans, ok := eb.subscribers[event.Name]; ok {
		for _, ch := range chans {
			go func(ch chan Event) {
				ch <- event
			}(ch)
		}
	}
}
