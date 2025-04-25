package client

import (
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/mohamidsaiid/uniclipboard/internal/clipboard"
)

type Client struct {
	conn           *websocket.Conn
	logger         *log.Logger
	closeConn      chan bool
	uniClipboard   chan []byte
	localClipboard chan []byte
}

func NewClient(URL url.URL) *Client {
	conn, err := newWebsocketConn(URL)
	if err != nil {
		panic(err)
	}

	return &Client{
		conn:         conn,
		logger:       log.New(os.Stdout, "", log.Ldate|log.Ltime),
		closeConn:    make(chan bool),
		uniClipboard: make(chan []byte),
	}
}

func (cl *Client) StartClient() error {
	go cl.reciveMessagesHandler()
	go clipboard.WatcheHandler(cl.localClipboard)

	for {
		select {
		case data := <-cl.localClipboard:
			cl.sendMessage(data)
		case data := <-cl.uniClipboard:
			clipboard.WriteHandler(data)
		case <-cl.closeConn:
			return errors.New("closing connection")
		}
	}
}
