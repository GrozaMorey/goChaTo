package client

import (
	"github.com/gorilla/websocket"
)

func main() {
	websocket.DefaultDialer.Dial("ws://127.0.0.1:8000/ws", nil)
}
