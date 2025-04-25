package uniclipboard

import (
	"context"
	"log"
	"time"

	"golang.design/x/clipboard"
)

type Message struct {
	Type clipboard.Format
	Data []byte
}

type UniClipboard struct {
	LocalClipboard Message
	UniClipboard   Message
	// the uniclipboard has a timeout
	TemporaryClipboardTimeout time.Duration
	// to indicate there is a new data written to the local clipboard
	NewDataWrittenLocaly chan struct{}
}

func init() {
	err := clipboard.Init()
	if err != nil {
		log.Fatalln(err)
	}
}

func (uc *UniClipboard) watchTextHandler() {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtText)

		uc.LocalClipboard = Message{
			Type: clipboard.FmtText,
		} 

		uc.LocalClipboard.Data = <-changed
		uc.NewDataWrittenLocaly <- struct{}{}
	}
}

func (uc *UniClipboard) watchImageHandler() {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtImage)

		uc.LocalClipboard = Message{
			Type: clipboard.FmtImage,
		}
		
		uc.LocalClipboard.Data = <-changed
		uc.NewDataWrittenLocaly <- struct{}{}
	}
}

func (uc *UniClipboard) WatchHandler() {
	go uc.watchImageHandler()
	go uc.watchTextHandler()
}

func (uc *UniClipboard) writeHandler(data Message) {
	clipboard.Write(data.Type, data.Data)
}

func (uc *UniClipboard) readHanlder(messageType clipboard.Format) Message {
	data := clipboard.Read(messageType)
	return Message{
		Type: messageType,
		Data: data,
	}
}

func (uc *UniClipboard) WriteTemporaryHanlder() {
	latestClipboardData := uc.readHanlder(uc.UniClipboard.Type)
	uc.writeHandler(uc.UniClipboard)
	time.Sleep(uc.TemporaryClipboardTimeout)
	uc.writeHandler(latestClipboardData)
}
