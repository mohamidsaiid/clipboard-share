package app

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
	"github.com/mohamidsaiid/uniclipboard/internal/clipboard"
	"github.com/mohamidsaiid/uniclipboard/internal/discovery"
	"github.com/mohamidsaiid/uniclipboard/internal/server"
)

type Application struct {
	conn *websocket.Conn
	logger *log.Logger
	closeConn chan bool
	uniClipboard chan []byte
	localClipboard chan []byte
}

func NewApplication(URL url.URL) *Application {
	conn, err := newWebsocketConn(URL)
	if err  != nil {
		panic(err)
	}

	return &Application{
		conn: conn,
		logger : log.New(os.Stdout, "", log.Ldate|log.Ltime),
		closeConn: make(chan bool),
		uniClipboard: make(chan []byte),
	}
}

func App(baseURL string, port string) error {
start:
	URL, err := discovery.ValidServer(baseURL, port, 2, 254)
	log.Println(URL.String())
	if err != nil {
		log.Println(err)
		srvr := server.NewServer(port)
		go srvr.Start()
		URL = &url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1%s",port), Path:"/clipboard"}
	}
	app := NewApplication(*URL)
	
	log.Println(client(app))
	goto start
}

func client(app *Application) error {
	go app.recivedWebsocketMessagesHandler(app.conn)
	go clipboard.WatcheHandler(app.localClipboard)

	for {
		select {
		case data := <- app.localClipboard:
			app.SendMessage(data)
		case data := <-app.uniClipboard:
			clipboard.WriteHandler(data)
		case <- app.closeConn:
			return errors.New("closing connection")
		}
	}
}
