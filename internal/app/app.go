package app

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/mohamidsaiid/uniclipboard/internal/client"
	"github.com/mohamidsaiid/uniclipboard/internal/discovery"
	"github.com/mohamidsaiid/uniclipboard/internal/server"
)

func StartApp(baseURL string, port string) error {
start:
	log.Println("Starting application...")

	log.Println("discovering valid server...")
	URL, err := discovery.ValidServer(baseURL, port, 2, 254)
	log.Println(URL.String())

	if err != nil {
		log.Println(err)
		srvr := server.NewServer(port)
		go srvr.Start()
		URL = url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1%s",port), Path:"/clipboard"}
	}

	time.Sleep(2 * time.Second)
	log.Println("Connecting to server...")

	cl := client.NewClient(URL)

	log.Println("Connected to server")
	log.Println(cl.StartClient())
	goto start
}
