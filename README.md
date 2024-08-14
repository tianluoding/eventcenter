# eventcenter

Golang implementation of eventcenter based on eventbus and websocket.

* eventbus
* websocket server

## Usage

```go
eb := eventbus.NewEventBus()
center := eventcenter.NewEventCenter(eb)

/**
* Encapsulate handle function on your need, eg. func(echo.Context), func(*gin.Context)... 
*/
http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    center.HandleWebSocket(w, r)
})
http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
    center.Publish(eventbus.Event{Name: "test", Data: "hello world"})
})

log.Fatal(http.ListenAndServe(":8080", nil))
```