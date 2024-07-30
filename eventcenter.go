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

// TODO: center管理connection
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

func (center *EventCenter) Publish(event eventbus.Event) {
	center.eb.Publish(event)
}

func (center *EventCenter) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	msg := &Message{}
	eventCh := make(chan eventbus.Event)
	for {
		err := conn.ReadJSON(msg)
		if err != nil {
			log.Printf("read json err: %v", err)
			close(eventCh)
			break
		}
		if msg.MsgType == "subscription" {
			center.eb.Subscribe(msg.ID, msg.Name, eventCh)
			go func() {
				for {
					select {
					case event, ok := <-eventCh:
						if !ok {
							log.Printf("channel closed, %v unsubscribe topic: %v", msg.ID, msg.Name)
							center.eb.Unsubscribe(msg.ID, msg.Name)
							return
						}
						if err := conn.WriteJSON(event); err != nil {
							log.Printf("write error: %v, %v unsubscribe topic: %v", err, msg.ID, msg.Name)
							close(eventCh)
							center.eb.Unsubscribe(msg.ID, msg.Name)
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
