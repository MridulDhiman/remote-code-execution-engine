package ws

import (
	"net/url"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	addr string
	url  url.URL
	Conn *websocket.Conn
}

func  NewWsClient(addr string) *WsClient {
	return &WsClient{
		addr: addr, 
		url: url.URL{Scheme: "ws", Host: addr, Path: "/echo"},
	}
}

func (w *WsClient) InitWsConnection() (error) {
	wsConn, _, err:= websocket.DefaultDialer.Dial(w.url.String(), nil)
	if err != nil {
		return err
	}

	w.Conn = wsConn
return nil
}

func (w* WsClient) SendMsg (msg string) error {
return w.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
}



