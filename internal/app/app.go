package app

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/mohamidsaiid/uniclipboard/internal/client"
	uniclipboard "github.com/mohamidsaiid/uniclipboard/internal/clipboard"
	"github.com/mohamidsaiid/uniclipboard/internal/discovery"
	"github.com/mohamidsaiid/uniclipboard/internal/server"
)

func StartApp(baseURL string, port string) error {
start:
	log.Println("Starting application...")

	clipboard, err := uniclipboard.NewClipboard(time.Minute * 15)
	if err != nil {
		return err
	}

	log.Println("discovering valid server...")
	link, err := discovery.ValidServer(baseURL, port, "/api/v1/healthcheck", 2, 254)

	if err != nil {
		log.Println(err)
		srvr := server.NewServer(port, clipboard)
		go srvr.Start()
		link = url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1%s", port), Path: "/api/v1/clipboard"}
	}

	time.Sleep(2 * time.Second)
	log.Println("Connecting to server...")

	link.Scheme = "ws"
	link.Path = "/api/v1/clipboard"
	cl, err := client.NewClient(link, clipboard)
	if err != nil {
		return err
	}

	log.Println("Connected to server")
	log.Println(cl.StartClient())
	goto start
}
