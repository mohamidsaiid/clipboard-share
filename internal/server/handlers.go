package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mohamidsaiid/uniclipboard/internal/jsonParser"
)

type Message struct {
	sender *websocket.Conn
	data []byte
}

var upgrader = websocket.Upgrader{}
var broadcast = make(chan Message)
var mutex = &sync.Mutex{}



func (srvr *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]bool{
		"ok": true,
	}

	err := jsonParser.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		srvr.serverErrorResponse(w, r, err)
	}

}

func (srvr *Server) clipboardHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		srvr.logError(r, err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	srvr.Clients[conn] = true
	mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			srvr.logError(r, err)
			mutex.Lock()
			delete(srvr.Clients, conn)
			mutex.Unlock()
			break
		}

		broadcast <- Message{sender: conn, data: message}
	}
}

func (srvr *Server) handleMessages() {
	for {
		message := <-broadcast

		mutex.Lock()
		for client := range srvr.Clients {
			if message.sender == client {
				continue
			}

			err := client.WriteMessage(websocket.TextMessage, message.data)
			if err != nil {
				client.Close()
				delete(srvr.Clients, client)
			}
		}
		mutex.Unlock()
	}
}
