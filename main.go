package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type msg struct {
	Num int
}

const port = 3000

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", websocketHandler)

	fmt.Printf("Server started on port %d\n", port)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not permitted", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", http.StatusBadRequest)
	}

	go echo(conn)
}

func echo(conn *websocket.Conn) {
	for {
		message := msg{}

		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Println("JSON read error", err)
		}

		fmt.Printf("Data received: %#v\n", message)
	}
}
