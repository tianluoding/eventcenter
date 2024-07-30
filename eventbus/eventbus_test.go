package eventbus

import (
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/smartystreets/goconvey/convey"
)

// TestNewEventBus checks the creation of a new EventBus instance.
func TestNewEventBus(t *testing.T) {
	convey.Convey("Given I create a new EventBus", t, func() {
		bus := NewEventBus()

		convey.So(bus, convey.ShouldNotBeNil)
		convey.So(bus.subscribers, convey.ShouldBeEmpty)
	})
}

// TestSubscribe tests the Subscribe method of EventBus.
func TestSubscribe(t *testing.T) {
	g := NewGomegaWithT(t)
	bus := NewEventBus()

	testCh := make(chan Event)
	eventName := "testEvent"
	id := "testID"

	bus.Subscribe(id, eventName, testCh)

	g.Eventually(func() bool {
		_, ok := bus.subscribers[eventName][id]
		return ok
	}).Should(BeTrue())
}

// TestUnsubscribe tests the Unsubscribe method of EventBus.
func TestUnsubscribe(t *testing.T) {
	g := NewGomegaWithT(t)
	bus := NewEventBus()

	testCh := make(chan Event)
	eventName := "testEvent"
	id := "testID"

	bus.Subscribe(id, eventName, testCh)
	bus.Unsubscribe(id, eventName)

	g.Eventually(func() bool {
		_, ok := bus.subscribers[eventName][id]
		return !ok
	}).Should(BeTrue())
}

// TestPublish tests the Publish method of EventBus.
func TestPublish(t *testing.T) {
	g := NewGomegaWithT(t)
	bus := NewEventBus()

	testCh := make(chan Event)
	eventName := "testEvent"
	id := "testID"

	bus.Subscribe(id, eventName, testCh)

	event := Event{Name: eventName, Data: "testData"}
	bus.Publish(event)

	select {
	case receivedEvent := <-testCh:
		g.Expect(receivedEvent).To(Equal(event))
	case <-time.After(time.Second):
		t.Errorf("Timed out waiting for event to be published")
	}
}

func TestPublishConcurrent(t *testing.T) {
	RegisterTestingT(t)

	bus := NewEventBus()
	testCh := make(chan Event, 10)
	eventName := "testEvent"
	id := "testID"

	bus.Subscribe(id, eventName, testCh)

	// 创建多个事件并同时发布
	events := []Event{
		{Name: eventName, Data: "testData1"},
		{Name: eventName, Data: "testData2"},
		{Name: eventName, Data: "testData3"},
	}

	var wg sync.WaitGroup
	wg.Add(len(events))

	for _, event := range events {
		go func(event Event) {
			defer wg.Done()
			bus.Publish(event)
		}(event)
	}

	// 等待所有事件发布完成
	wg.Wait()

	// 使用Eventually检查所有事件最终都被接收到
	Eventually(testCh, time.Second*5, time.Millisecond*100).Should(HaveLen(len(events)))

	// 检查每个事件是否正确接收到
	receivedEvents := make([]Event, 0, len(events))
	for i := 0; i < len(events); i++ {
		receivedEvent := <-testCh
		receivedEvents = append(receivedEvents, receivedEvent)
	}

	// 验证接收到的事件数据
	Expect(receivedEvents).To(ConsistOf(events))
}
