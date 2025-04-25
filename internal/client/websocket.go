package client

import (
	"net/url"

	"github.com/gorilla/websocket"
	uniclipboard "github.com/mohamidsaiid/uniclipboard/internal/clipboard"
	"golang.design/x/clipboard"
)

func newWebsocketConn(url url.URL) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cl *Client) sendMessage() error {
	var messageType int
	if cl.clipboard.UniClipboard.Type == clipboard.FmtText {
		messageType = websocket.TextMessage
	} else {
		messageType = websocket.BinaryMessage
	}
	
	err := cl.conn.WriteMessage(messageType, cl.clipboard.UniClipboard.Data)
	if err != nil {
		return err
	}

	return nil
}

func (cl *Client) receiveMessage() *uniclipboard.Message {
	messageType, message, err := cl.conn.ReadMessage()
	if err != nil {
		cl.logger.Println("read: ", err)
		cl.close()
		return nil
	}
	cl.newWrittenDataUni <- struct{}{}	
	return &uniclipboard.Message{
		Type: clipboard.Format(messageType),
		Data: message,
	}
}

func (cl *Client) close() error {
	err := cl.conn.Close()
	cl.closeConn <- true
	if err != nil {
		return err
	}
	return nil
}

func (cl *Client) reciveWebsocketMessagesHandler() {
	for {
		cl.clipboard.UniClipboard = cl.receiveMessage()
	}
}

