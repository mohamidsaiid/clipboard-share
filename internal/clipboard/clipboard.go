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
	UniClipboard   *Message
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

		uc.UniClipboard = &Message{
			Type: clipboard.FmtText,
		} 

		uc.UniClipboard.Data = <-changed
		uc.NewDataWrittenLocaly <- struct{}{}
	}
}

func (uc *UniClipboard) watchImageHandler() {
	for {
		changed := clipboard.Watch(context.Background(), clipboard.FmtImage)

		uc.UniClipboard = &Message{
			Type: clipboard.FmtImage,
		}
		
		uc.UniClipboard.Data = <-changed
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

func (uc *UniClipboard) ReadHanlder(messageType clipboard.Format) Message {
	data := clipboard.Read(messageType)
	return Message{
		Type: messageType,
		Data: data,
	}
}

func (uc *UniClipboard) WriteTemporaryHanlder() {
	if uc.UniClipboard == nil {
        log.Fatal("UniClipboard instance is nil")
        return
    }
    log.Printf("UniClipboard state: %+v", uc.UniClipboard)
	// save the latest clipboard data
	latestClipboardData := uc.ReadHanlder(uc.UniClipboard.Type)
	// write the new uniclipboard data	
	uc.writeHandler(*uc.UniClipboard)
	// wait for the specified time till the uniclipboard is vanished
	time.Sleep(uc.TemporaryClipboardTimeout)
	// rewrite the old localclipboard data
	uc.writeHandler(latestClipboardData)
	// remove the data from the uniclipboard
	uc.UniClipboard = nil
}
