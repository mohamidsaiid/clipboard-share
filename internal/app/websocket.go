package app
import (
	"net/url"

	"github.com/gorilla/websocket"
)

func (app *Application) SendMessage(message []byte) error {
	err := app.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) ReceiveMessage() ([]byte, error) {
	_, message, err := app.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (app *Application) Close() error {
	err := app.conn.Close()
	app.closeConn <- true
	if err != nil {
		return err
	}
	return nil
}

func newWebsocketConn(url url.URL) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (app *Application) recivedWebsocketMessagesHandler(conn *websocket.Conn) {
	for {
		_, message, err := app.conn.ReadMessage()
		if err != nil {
			app.logger.Println("read: ", err)
			app.Close()
			return
		}
		app.uniClipboard <- message
	}
}