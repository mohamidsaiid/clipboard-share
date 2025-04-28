package server

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	uniclipboard "github.com/mohamidsaiid/uniclipboard/internal/clipboard"
	"github.com/mohamidsaiid/uniclipboard/internal/jsonParser"
	"golang.design/x/clipboard"
)

type Message struct {
	sender *websocket.Conn
	data   uniclipboard.Message
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
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			srvr.logError(r, err)
			mutex.Lock()
			delete(srvr.Clients, conn)
			mutex.Unlock()
			break
		}

		if messageType == websocket.BinaryMessage {
			broadcast <- Message{sender: conn, data: uniclipboard.Message{
				Type: clipboard.FmtImage,
				Data: message,
			}}
			if messageType == websocket.TextMessage {
				broadcast <- Message{sender: conn, data: uniclipboard.Message{
					Type: clipboard.FmtText,
					Data: message,
				}}
			}
		}
	}
}

func (srvr *Server) handleMessages() {
	for {
		message := <-broadcast
		var messageType int
		if message.data.Type == clipboard.FmtText {
			messageType = websocket.TextMessage
		} else {
			messageType = websocket.BinaryMessage
		}

		mutex.Lock()
		for client := range srvr.Clients {
			if message.sender == client {
				continue
			}

			err := client.WriteMessage(messageType, message.data.Data)
			if err != nil {
				client.Close()
				delete(srvr.Clients, client)
			}
		}
		mutex.Unlock()
	}
}

func (srvr *Server) lastCopiedData(w http.ResponseWriter, r *http.Request) {
}
//, "test" : []any{"uniclip", srvr.clipboard}
//, "test" : []any{"origiclip", srvr.clipboard}