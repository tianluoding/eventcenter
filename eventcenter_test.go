package eventcenter

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tianluoding/eventcenter/eventbus"
)

func TestEventCenter(t *testing.T) {
	Convey("Given an EventCenter instance", t, func() {
		eb := &eventbus.EventBus{}
		center := NewEventCenter(eb)

		Convey("When publishing an event", func() {
			var called bool
			patches := gomonkey.ApplyMethod(eb, "Publish", func(eb *eventbus.EventBus, event eventbus.Event) {
				called = true
			})
			defer patches.Reset()

			center.Publish(eventbus.Event{})
			So(called, ShouldBeTrue)
		})

		Convey("When listening and writing to WebSocket", func() {
			fakeConn := &websocket.Conn{}

			eventCh := make(chan eventbus.Event)
			handleFinished := make(chan bool)
			msg := &Message{ID: "testID", MsgType: "subscription"}
			patches := gomonkey.ApplyMethodReturn(fakeConn, "WriteJSON", nil)
			defer patches.Reset()

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				center.listenAndWrite(eventCh, fakeConn, msg, handleFinished)
			}()

			// Simulate receiving an event
			eventCh <- eventbus.Event{}
			time.Sleep(10 * time.Millisecond) // Allow time for the goroutine to process the event

			close(eventCh)
			wg.Wait()
		})

		Convey("When listening and writing to WebSocket, handleFinished", func() {
			fakeConn := &websocket.Conn{}

			eventCh := make(chan eventbus.Event)
			handleFinished := make(chan bool)
			msg := &Message{ID: "testID", MsgType: "subscription"}
			patches := gomonkey.ApplyMethodReturn(fakeConn, "WriteJSON", nil)
			defer patches.Reset()

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				center.listenAndWrite(eventCh, fakeConn, msg, handleFinished)
			}()

			// Simulate receiving an event
			eventCh <- eventbus.Event{}
			time.Sleep(10 * time.Millisecond) // Allow time for the goroutine to process the event

			handleFinished <- true
			wg.Wait()
		})

		Convey("When listening and writing to WebSocket, ReadJSON err", func() {
			fakeConn := &websocket.Conn{}

			patches := gomonkey.
				ApplyMethodReturn(&upgrader, "Upgrade", fakeConn, nil).
				ApplyMethodReturn(fakeConn, "Close", nil).
				ApplyMethodReturn(fakeConn, "ReadJSON", fmt.Errorf("read error"))

			defer patches.Reset()

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				center.HandleWebSocket(nil, nil)
			}()

			wg.Wait()
		})
	})
}
