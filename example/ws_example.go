package main

import (
	"log"
	"net/http"

	"github.com/tianluoding/eventcenter"
	"github.com/tianluoding/eventcenter/eventbus"
)

func main() {
	eb := eventbus.NewEventBus()
	center := eventcenter.NewEventCenter(*eb)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		center.HandleWebSocket(w, r)
	})
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		eb.Publish(eventbus.Event{Name: "test", Data: "hello world"})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
