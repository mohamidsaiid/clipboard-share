package client

import (
	"errors"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mohamidsaiid/uniclipboard/internal/clipboard"
)

type Client struct {
	conn      *websocket.Conn
	logger    *log.Logger
	closeConn chan bool
	clipboard *uniclipboard.UniClipboard
	newWrittenDataUni chan struct{}
}

func NewClient(URL url.URL) (*Client, error) {
	conn, err := newWebsocketConn(URL)
	if err != nil {
		return nil, err
	}
	clipboard := &uniclipboard.UniClipboard{
		LocalClipboard: uniclipboard.Message{},
		UniClipboard: uniclipboard.Message{},
		TemporaryClipboardTimeout: time.Minute * 15,
		NewDataWrittenLocaly: make(chan struct{}),
	}

	return &Client{
		conn:      conn,
		logger:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		closeConn: make(chan bool),
		clipboard: clipboard,
		newWrittenDataUni: make(chan struct{}),
	}, nil
}

func (cl *Client) StartClient() error {
	go cl.reciveWebsocketMessagesHandler()
	go cl.clipboard.WatchHandler()

	for {
		select {
		case <-cl.clipboard.NewDataWrittenLocaly:
			cl.sendMessage()
		case <-cl.newWrittenDataUni:
			go cl.clipboard.WriteTemporaryHanlder()
		case <-cl.closeConn:
			return errors.New("closing connection")
		}
	}
}
