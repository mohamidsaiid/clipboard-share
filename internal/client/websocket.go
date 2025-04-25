package client

import (
	"net/url"

	"github.com/gorilla/websocket"
)

func newWebsocketConn(url url.URL) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cl *Client) sendMessage(message []byte) error {
	err := cl.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (cl *Client) receiveMessage() []byte {
	_, message, err := cl.conn.ReadMessage()
	if err != nil {
		cl.logger.Println("read: ", err)
		cl.close()
		return nil
	}
	return message
}

func (cl *Client) close() error {
	err := cl.conn.Close()
	cl.closeConn <- true
	if err != nil {
		return err
	}
	return nil
}

func (cl *Client) reciveMessagesHandler() {
	for {
		cl.uniClipboard <- cl.receiveMessage()
	}
}
