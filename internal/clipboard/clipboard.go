package uniclipboard

import (
	"context"
	"log"
	"sync"

	"github.com/mohamidsaiid/uniclipboard/internal/ADT"
	"golang.design/x/clipboard"
)

type Message struct {
	Type clipboard.Format
	Data []byte
}

type UniClipboard struct {
	UniClipboard Message
	// to indicate there is a new data written to the local clipboard
	NewDataWrittenLocaly ADT.Sig
	Mutex                *sync.Mutex
}

func NewClipboard(sig ADT.Sig) (*UniClipboard, error) {
	err := clipboard.Init()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &UniClipboard{
		UniClipboard:              Message{},
		NewDataWrittenLocaly:      sig,
		Mutex:                     &sync.Mutex{},
	}, nil
}

func (uc *UniClipboard) watchTextHandler() {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtText)

		data := <-changed
		log.Println("clipboard package new text data been written to the clipboard internally")
		log.Println(data)

		uc.Mutex.Lock()
		uc.UniClipboard.Data = data
		uc.UniClipboard.Type = clipboard.FmtText
		uc.NewDataWrittenLocaly <- struct{}{}
		uc.Mutex.Unlock()
	}
}

func (uc *UniClipboard) watchImageHandler() {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtImage)

		data := <-changed
		log.Println("clipboard package new image data been written to the clipboard internally")
		log.Println(data)

		uc.Mutex.Lock()
		uc.UniClipboard.Data = data
		uc.UniClipboard.Type = clipboard.FmtImage
		uc.NewDataWrittenLocaly <- struct{}{}
		uc.Mutex.Unlock()

	}
}

func (uc *UniClipboard) WatchHandler() {
	go uc.watchImageHandler()
	go uc.watchTextHandler()
}

func (uc *UniClipboard) WriteHandler(data Message) {
	log.Println("clipboard package new data is going to be written to the clipboard ", data)
	clipboard.Write(data.Type, data.Data)
}

func (uc *UniClipboard) ReadHanlder(messageType clipboard.Format) Message {
	data := clipboard.Read(messageType)
	log.Println("clipboard packag data been read from the internal clipboard")
	return Message{
		Type: messageType,
		Data: data,
	}
}