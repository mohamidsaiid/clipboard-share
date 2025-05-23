package client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"golang.design/x/clipboard"
)

func newWebsocketConn(url url.URL, sk string) (*websocket.Conn, error) {
	rqHeader := http.Header{}
	authKey := fmt.Sprint("Bearer ",sk)
	rqHeader.Add("Authorization", authKey)

	conn, _, err := websocket.DefaultDialer.Dial(url.String(), rqHeader)
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

func (cl *Client) receiveMessage() {
	for {
		messageType, message, err := cl.conn.ReadMessage()
		if err != nil {
			cl.logger.Println("read: ", err)
			cl.close()
			continue
		}

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
