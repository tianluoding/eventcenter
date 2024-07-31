package eventcenter

import (
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tianluoding/eventcenter/eventbus"
)

func TestEventCenter(t *testing.T) {
	Convey("Given an EventCenter instance", t, func() {
		eb := &eventbus.EventBus{}
		center := NewEventCenter(eb)

		Convey("When publishing an event", func() {
			patches := gomonkey.ApplyFunc(eb.Publish, func(event eventbus.Event) {})
			defer patches.Reset()

			center.Publish(eventbus.Event{})
		})

		Convey("When handling WebSocket connections", func() {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/", nil)

			patches := gomonkey.ApplyFunc(eb.Subscribe, func(id string, name string, ch chan eventbus.Event) {})
			patches.ApplyFunc(eb.Unsubscribe, func(id string, name string) {})
			defer patches.Reset()

			center.HandleWebSocket(w, r)
		})
	})
}
