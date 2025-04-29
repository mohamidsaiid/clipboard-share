package client

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"golang.design/x/clipboard"
)

func newWebsocketConn(url url.URL) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	log.Println("connected to the server")
	log.Println("url: ", url.String())
	return conn, nil
}

func (cl *Client) sendMessage() error {
	var messageType int

	if cl.clipboard.UniClipboard.Type == clipboard.FmtText {
		messageType = websocket.TextMessage
	} else {
		messageType = websocket.BinaryMessage
	}

	log.Println("sending message")
	log.Println(cl.clipboard.UniClipboard.Type, string(cl.clipboard.UniClipboard.Data))

	err := cl.conn.WriteMessage(messageType, cl.clipboard.UniClipboard.Data)
	if err != nil {
		return err
	}

	return nil
}

func (cl *Client) receiveMessage() {
	for {
		messageType, message, err := cl.conn.ReadMessage()
		if err != nil {
			cl.logger.Println("read: ", err)
			cl.close()
			continue
		}

		log.Println("client/websocket new received message")
		log.Println(messageType, string(message))

		cl.newWrittenDataUni <- struct{}{}

		if messageType == websocket.TextMessage {
			messageType = int(clipboard.FmtText)
		} else {
			messageType = int(clipboard.FmtImage)
		}

		cl.clipboard.Mutex.Lock()
		cl.clipboard.UniClipboard.Data = message
		cl.clipboard.UniClipboard.Type = clipboard.Format(messageType)
		cl.clipboard.Mutex.Unlock()

	}
}

func (cl *Client) close() error {
	err := cl.conn.Close()
	cl.closeConn <- struct{}{}
	if err != nil {
		return err
	}
	return nil
}
