package eventcenter

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tianluoding/eventcenter/eventbus"
)

type Message struct {
	ID      string `json:"id"`
	MsgType string `json:"type"`
	eventbus.Event
}

type EventCenter struct {
	eb eventbus.EventBus
}

func NewEventCenter(eb eventbus.EventBus) *EventCenter {
	return &EventCenter{eb: eb}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (center *EventCenter) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		msg := &Message{}
		err := conn.ReadJSON(msg)
		if err != nil {
			log.Printf("read json err: %v", err)
		}
		if msg.MsgType == "subscription" {
			eventCh := make(chan eventbus.Event)
			center.eb.Subscribe(msg.ID, msg.Name, eventCh)
			go func() {
				for {
					select {
					case event, ok := <-eventCh:
						if !ok {
							log.Println("channel closed")
							return
						}
						if err := conn.WriteJSON(event); err != nil {
							log.Println("Write error:", err)
							return
						}
					}
				}
			}()
		} else if msg.MsgType == "unsubscription" {
			center.eb.Unsubscribe(msg.ID, msg.Name)
			time.Sleep(time.Second * 1)
			break
		}
	}
}
