package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Uuid string `json:"uuid"`
	Ice  struct {
		Candidate        string `json:"candidate"`
		SdpMid           string `json:"sdpMid"`
		SdpMLineIndex    string `json:"sdpMLineIndex"`
		Protocol         string `json:"protocol"`
		Foundation       string `json:"foundation"`
		Priority         string `json:"priority"`
		Component        string `json:"component"`
		Port             string `json:"port"`
		Address          string `json:"address"`
		Type             string `json:"type"`
		TcpType          string `json:"tcpType"`
		UsernameFragment string `json:"usernameFragment"`
	} `json:"ice"`
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func main() {
	port := getPort()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", websocketHandler)

	fmt.Printf("Server started on port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "https://"+r.Host {
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
		message := Message{}

		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Println("JSON read error", err)
		}

		fmt.Printf("Data received: %#v\n", message)
		broadcastMessagesToClients()
	}
}

func broadcastMessagesToClients() {
	for {
		message := <-broadcast

		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error occured writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func getPort() string {
	value := os.Getenv("PORT")
	if len(value) == 0 {
		return "3000"
	}
	return value
}
