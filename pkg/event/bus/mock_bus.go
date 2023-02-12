package bus

import (
	"sync"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

type EventBusMock struct {
	events sync.Map
}

func NewEventBusMock() *EventBusMock {
	return &EventBusMock{}
}

func (b *EventBusMock) Publish(event testkube.Event) error {
	b.events.Range(func(key, e interface{}) bool {
		e.(chan testkube.Event) <- event
		return true
	})
	return nil
}
func (b *EventBusMock) Subscribe(queue string, handler Handler) error {

	ch := make(chan testkube.Event)
	b.events.Store(queue, ch)

	go func() {
		for e := range ch {
			handler(e)
		}
	}()
	return nil
}

func (b *EventBusMock) Unsubscribe(queue string) error {
	return nil

}
func (b *EventBusMock) Close() error {
	b.events.Range(func(key, e interface{}) bool {
		b.events.Delete(key)
		return true
	})
	return nil
}
