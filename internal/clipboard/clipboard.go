package uniclipboard

import (
	"context"
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
	clipboard.Write(data.Type, data.Data)
}

func (uc *UniClipboard) ReadHanlder(messageType clipboard.Format) Message {
	data := clipboard.Read(messageType)

	return Message{
		Type: messageType,
		Data: data,
	}
}