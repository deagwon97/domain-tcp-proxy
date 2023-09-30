package server

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

// echo websocket server
// send same message to client when receive the message from client
func echoServer(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}

	defer c.Close()
	for {
		// recieve
		mt, clientReadMessage, err := c.ReadMessage()
		if err != nil {
			break
		}
		// send
		err = c.WriteMessage(mt, clientReadMessage)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}

func RunServer(waitWebSocketServerReady *sync.WaitGroup, host string) {
	defer waitWebSocketServerReady.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", echoServer)
	listener, _ := net.Listen("tcp", host)
	waitWebSocketServerReady.Done()
	// fmt.Println("Starting server ", host)
	http.Serve(listener, mux)
}
