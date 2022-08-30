package doljara

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Room)

func ServWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade: ", err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var message Room
		if err := conn.ReadJSON(&message); err != nil {
			log.Println("ReadMessage: ", err)
			return
		}

		broadcast <- message
	}
}

func HandleMessage() {
	for {
		message := <-broadcast

		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				log.Println("WriteMessage: ", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
