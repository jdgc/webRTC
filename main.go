package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

type msg struct {
	Num int
}

func main() {
	port := os.Getenv("PORT")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", websocketHandler)

	fmt.Printf("Server started on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		log.Fatalf("ORIGIN not permitted: %s", r.Header.Get("Origin"))
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
