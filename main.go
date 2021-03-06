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
		SdpMid           string `json:"sdpMid,omitempty"`
		SdpMLineIndex    int    `json:"sdpMLineIndex,omitempty"`
		Foundation       string `json:"foundation,omitempty"`
		Component        string `json:"component,omitempty"`
		Priority         int    `json:"priority,omitempty"`
		Address          string `json:"address,omitempty"`
		Protocol         string `json:"protocol,omitempty"`
		Port             int    `json:"port,omitempty"`
		Type             string `json:"type,omitempty"`
		TcpType          string `json:"tcpType,omitempty"`
		RelatedAddress   string `json:"relatedAddress,omitempty"`
		RelatedPort      int    `json:"relatedPort,omitempty"`
		UsernameFragment string `json:"usernameFragment,omitempty"`
	} `json:"ice,omitempty"`
	Sdp struct {
		Type string `json:"type"`
		Sdp  string `json:"sdp"`
	} `json:"sdp,omitempty"`
}

var clients = make(map[*websocket.Conn]bool)

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
		log.Printf("ORIGIN not permitted: %s", r.Header.Get("Origin"))
		http.Error(w, "Origin not permitted", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", http.StatusBadRequest)
	}

	clients[conn] = true
	go echo(conn)
}

func echo(conn *websocket.Conn) {
	for {
		message := Message{}

		err := conn.ReadJSON(&message)
		if err != nil {
			fmt.Println("JSON read error", err)
		}

		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("error occured writing message to client: %v", err)
				conn.Close()
			}
			fmt.Printf("sent broadcast: %v", message)
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
