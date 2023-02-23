package ws

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WsWrapper struct {
	Conn *websocket.Conn
}

func (wsw *WsWrapper) Write(p []byte) error {
	err := wsw.Conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return err
	}

	return nil
}

func (wsw *WsWrapper) Read() ([]byte, error) {
	for {
		msgType, content, err := wsw.Conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return nil, err
		}

		if msgType != websocket.BinaryMessage {
			err := errors.New("websocket read message type error.")
			log.Error(err)
			log.Info(string(content))
			return nil, err
		}

		return content, err
	}
}
