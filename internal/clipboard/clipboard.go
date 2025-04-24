package clipboard

import (
	"context"
	"log"

	"golang.design/x/clipboard"
)

func init() {
	err := clipboard.Init()
	if err != nil {
		log.Fatalln(err)
	}
}

func WatcheHandler(localClipboard chan []byte) {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtText)
		localClipboard <- <-changed
	}
}

func WriteHandler(data []byte) {
	clipboard.Write(clipboard.FmtText, data)
}