package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("received: %s", msg)

		err = conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func serveHttp(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("{\"ping\": \"pong\"}"))
}

func main() {
	http.HandleFunc("/ws", serveWs)
	http.HandleFunc("/ping", serveHttp)

	srv := &http.Server{
		Handler:      http.DefaultServeMux,
		Addr:         "127.0.0.1:3765",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
