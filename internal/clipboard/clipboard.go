package uniclipboard

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mohamidsaiid/uniclipboard/internal/ADT"
	"golang.design/x/clipboard"
)

type Message struct {
	Type clipboard.Format
	Data []byte
}

type UniClipboard struct {
	UniClipboard Message
	// the uniclipboard has a timeout
	TemporaryClipboardTimeout time.Duration
	// to indicate there is a new data written to the local clipboard
	NewDataWrittenLocaly ADT.Sig
	Mutex                *sync.Mutex
}

func NewClipboard(timeOut time.Duration, sig ADT.Sig) (*UniClipboard, error) {
	err := clipboard.Init()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &UniClipboard{
		UniClipboard:              Message{},
		TemporaryClipboardTimeout: timeOut,
		NewDataWrittenLocaly:      sig,
		Mutex:                     &sync.Mutex{},
	}, nil
}

func (uc *UniClipboard) watchTextHandler() {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtText)

		data := <-changed
		log.Println("clipboard package new data been written to the clipboard internally")

		uc.Mutex.Lock()
		uc.UniClipboard.Data = data
		uc.UniClipboard.Type = clipboard.FmtImage
		uc.NewDataWrittenLocaly <- struct{}{}
		uc.Mutex.Unlock()
	}
}

func (uc *UniClipboard) watchImageHandler() {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtImage)

		data := <-changed
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

func (uc *UniClipboard) writeHandler(data Message) {
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

func (uc *UniClipboard) WriteTemporaryHanlder() {
	/*if uc.UniClipboard == {
		log.Fatal("UniClipboard instance is nil")
		return
	}
	*/
	uc.Mutex.Lock()
	// save the latest clipboard data
	latestClipboardData := uc.ReadHanlder(uc.UniClipboard.Type)
	// write the new uniclipboard data
	uc.writeHandler(uc.UniClipboard)
	// wait for the specified time till the uniclipboard is vanished
	time.Sleep(uc.TemporaryClipboardTimeout)
	// rewrite the old localclipboard data
	uc.writeHandler(latestClipboardData)
	// remove the data from the uniclipboard
	uc.UniClipboard = Message{}
	uc.Mutex.Unlock()
}
