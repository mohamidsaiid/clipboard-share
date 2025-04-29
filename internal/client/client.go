package client

import (
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/mohamidsaiid/uniclipboard/internal/ADT"
	"github.com/mohamidsaiid/uniclipboard/internal/clipboard"
)

type Client struct {
	conn              *websocket.Conn
	logger            *log.Logger
	closeConn         ADT.Sig
	clipboard         *uniclipboard.UniClipboard
	newWrittenDataUni ADT.Sig
}

func NewClient(URL url.URL, clipboard *uniclipboard.UniClipboard) (*Client, error) {
	URL.Scheme = "ws"
	URL.Path = "/api/v1/clipboard"
	conn, err := newWebsocketConn(URL)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:              conn,
		logger:            log.New(os.Stdout, "", log.Ldate|log.Ltime),
		closeConn:         make(ADT.Sig),
		clipboard:         clipboard,
		newWrittenDataUni: make(ADT.Sig),
	}, nil
}

func (cl *Client) StartClient() error {
	go cl.receiveMessage()
	go cl.clipboard.WatchHandler()

	for {
		select {
		case <-cl.clipboard.NewDataWrittenLocaly:
			log.Println("new written data localy signal recieved")
			log.Println(string(cl.clipboard.UniClipboard.Data))
			cl.sendMessage()
		case <-cl.newWrittenDataUni:
			log.Println("new written data uni signal recieved")
			log.Println(string(cl.clipboard.UniClipboard.Data))
			cl.clipboard.Mutex.Lock()
			cl.clipboard.WriteHandler(cl.clipboard.UniClipboard)
			cl.clipboard.Mutex.Unlock()
		case <-cl.closeConn:
			log.Println("the conncetion is closed signal recieved")
			return errors.New("closing connection")
		}
	}
}
