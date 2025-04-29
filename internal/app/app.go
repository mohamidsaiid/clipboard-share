package app

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/mohamidsaiid/uniclipboard/internal/ADT"
	"github.com/mohamidsaiid/uniclipboard/internal/client"
	uniclipboard "github.com/mohamidsaiid/uniclipboard/internal/clipboard"
	"github.com/mohamidsaiid/uniclipboard/internal/discovery"
	"github.com/mohamidsaiid/uniclipboard/internal/models"
	"github.com/mohamidsaiid/uniclipboard/internal/secretkey"
	"github.com/mohamidsaiid/uniclipboard/internal/server"
)

func StartApp(baseURL string, port string, secretKeyPort string, originalSecretKey string) error {
start:
	log.Println("Starting application...")

	clipboard, err := uniclipboard.NewClipboard(make(ADT.Sig))
	if err != nil {
		return err
	}

	userModel, err := models.InitateDatabase("users.db")
	if err != nil {
		return err
	}

	go secretkey.StartSecertKeyWebServer(secretKeyPort, userModel)	
	fmt.Printf("\n\nplease visit \"localhost%s/secretkey\" \nand provide a the new secretkey to be used all over your devices\n\n", secretKeyPort)

	sk, exists := userModel.Get()
	if !exists {
		userModel.Update(originalSecretKey)
		sk.SecretKey = originalSecretKey
	}
	log.Println("discovering valid server...")
	link, err := discovery.ValidServer(baseURL, port, "/api/v1/healthcheck", 2, 254)

	if err != nil {
		log.Println(err)
		srvr := server.NewServer(port, clipboard, userModel)
		go srvr.Start()
		link = url.URL{Scheme: "ws", Host: fmt.Sprintf("127.0.0.1%s", port), Path: "/api/v1/clipboard"}
	}

	time.Sleep(2 * time.Second)
	log.Println("Connecting to server...")

	cl, err := client.NewClient(link, clipboard, sk.SecretKey)
	if err != nil {
		return err
	}

	log.Println("Connected to server")
	log.Println(cl.StartClient())
	goto start
}
