package main

import (
	"flag"
	"goChato/scr/server"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8000", "http server address")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешает все запросы
	},
}

func statusOk(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")

}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Print("cant read da message")
		}

		if string(msg) == "END" {
			break
		}
		log.Print(string(msg))
	}
}

func main() {
	flag.Parse()

	server, err := server.New()
	if err != nil {
		return
	}

	http.HandleFunc("/", statusOk)
	http.HandleFunc("/ws", handleConnection)
	http.HandleFunc("/reg", server.HandleRegister)
	http.HandleFunc("/log", server.HandleLogin)

	httpServer := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Println("listening on addr", *addr)
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

// root: TreeNode{val: 4, left: TreeNode{val: 2, left: TreeNode{val: 1, left: None, right: None}, right: TreeNode{val: 3, left: None, right: None}}, right: TreeNode{val: 7, left: TreeNode{val: 6, left: None, right: None}, right: TreeNode{val: 9, left: None, right: None}}}
// left: TreeNode{val: 2, left: TreeNode{val: 1, left: None, right: None}, right: TreeNode{val: 3, left: None, right: None}}
// right: TreeNode{val: 7, left: TreeNode{val: 6, left: None, right: None}, right: TreeNode{val: 9, left: None, right: None}}
// root: TreeNode{val: 2, left: TreeNode{val: 1, left: None, right: None}, right: TreeNode{val: 3, left: None, right: None}}
// left: TreeNode{val: 1, left: None, right: None}
// right: TreeNode{val: 3, left: None, right: None}
// root: None
