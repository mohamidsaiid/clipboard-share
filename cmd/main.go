package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamidsaiid/uniclipboard/internal/discovery"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		return
	}
	serverURL, err := discovery.ValidServer(os.Getenv("BASE_URL"), os.Getenv("PORT"), 2, 254)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(serverURL)
	log.Println(serverURL.String())
}