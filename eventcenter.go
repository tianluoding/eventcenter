package eventcenter

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tianluoding/eventcenter/eventbus"
)

type EventCenter struct {
	eb eventbus.EventBus
}

type Message struct {
	id string
	eventbus.Event
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (center *EventCenter) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		msg := Message{}
		err := conn.ReadJSON(msg)
		if err != nil {
			log.Printf("read json %v", err)
		}
		if msg.Name == "subscription" {
			eventCh := make(chan eventbus.Event)
			center.eb.Subscribe(msg.id, msg.Name, eventCh)
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
		} else if msg.Name == "unsubscription" {
			center.eb.Unsubscribe(msg.id, msg.Name)
		}
	}
}
