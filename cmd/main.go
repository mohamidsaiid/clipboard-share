package main

import (
	"log"
	"os"

	"github.com/mohamidsaiid/uniclipboard/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return
	}
	
	log.Println(app.StartApp(os.Getenv("BASE_URL"), os.Getenv("PORT")))
}