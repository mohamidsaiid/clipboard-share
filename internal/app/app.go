package app

import (
	"fmt"
	"log"
	"net/url"
	"time"

	uniclipboard "github.com/mohamidsaiid/uniclipboard/internal/clipboard"
	"github.com/mohamidsaiid/uniclipboard/internal/client"
_	"github.com/mohamidsaiid/uniclipboard/internal/discovery"
	"github.com/mohamidsaiid/uniclipboard/internal/server"
)

func StartApp(baseURL string, port string) error {
start:
	log.Println("Starting application...")

	log.Println("discovering valid server...")
	//URL, err := discovery.ValidServer(baseURL, port, 2, 254)
	//log.Println(URL.String())

	clipboard := &uniclipboard.UniClipboard{
		UniClipboard: nil,
		TemporaryClipboardTimeout: time.Minute * 15,
		NewDataWrittenLocaly: make(chan struct{}),
	}
	//if err != nil {
		//log.Println(err)
		srvr := server.NewServer(port, clipboard)
		go srvr.Start()
		URL := url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1%s",port), Path:"/clipboard"}
	//}

	time.Sleep(2 * time.Second)
	log.Println("Connecting to server...")
	
	cl, err := client.NewClient(URL, clipboard)
	if err != nil {
		return err
	}

	log.Println("Connected to server")
	log.Println(cl.StartClient())
	goto start
}
