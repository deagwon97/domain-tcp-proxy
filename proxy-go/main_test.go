package main_test

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"testing"
	"time"
	"tunnel/lib"

	"github.com/gorilla/websocket"
)

const (
	NUM_OF_SERVER = 3 // fixed
	NUM_REPEAT    = 10000
	NUM_OF_CLIENT = 100
	// PACKET_SIZE         = 1024 // minimum 128
	TUNNEL_HOST_POSTFIX = ".hi:9980"
)

// var data = make([]byte, PACKET_SIZE)

var upgrader = websocket.Upgrader{} // use default options

// test http based proxy server performance
func TestWebScoket(t *testing.T) {
	waitWebSocketServerReady := &sync.WaitGroup{}
	fmt.Println("Starting server ")
	waitWebSocketServerReady.Add(NUM_OF_SERVER)
	appHost1 := "0.0.0.0:8080"
	go runServer(waitWebSocketServerReady, appHost1)
	appHost2 := "0.0.0.0:8081"
	go runServer(waitWebSocketServerReady, appHost2)
	appHost3 := "0.0.0.0:8082"
	go runServer(waitWebSocketServerReady, appHost3)
	waitWebSocketServerReady.Wait()
	waitAllClientsEnd := &sync.WaitGroup{}
	fmt.Println("Starting client ")
	t.Log("Starting client ")
	waitAllClientsEnd.Add(NUM_OF_CLIENT * NUM_OF_SERVER)
	for i := 0; i < NUM_OF_CLIENT; i++ {
		go runClient(waitAllClientsEnd, 1, appHost1, t)
		go runClient(waitAllClientsEnd, 2, appHost2, t)
		go runClient(waitAllClientsEnd, 3, appHost3, t)
	}
	waitAllClientsEnd.Wait()
}

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
			fmt.Println("close server!")
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

func runServer(waitWebSocketServerReady *sync.WaitGroup, host string) {
	defer waitWebSocketServerReady.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", echoServer)
	listener, _ := net.Listen("tcp", host)
	waitWebSocketServerReady.Done()
	http.Serve(listener, mux)
}

func runClient(waitClientEnd *sync.WaitGroup, clientId int, appHost string, t *testing.T) {
	defer waitClientEnd.Done()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	subdomain, _ := lib.EncryptSubdomain(appHost)
	tunnelHost := subdomain + TUNNEL_HOST_POSTFIX
	fmt.Println("Connecting to " + tunnelHost)
	t.Log("Connecting to " + tunnelHost)
	u := url.URL{Scheme: "ws", Host: tunnelHost, Path: "/echo"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// map for check send and recieve same messages
	sendMap := make(map[string]bool)
	recvMap := make(map[string]bool)

	// 수신부
	done := make(chan struct{})
	waitSendRecieveMessageMap := &sync.WaitGroup{}
	waitSendRecieveMessageMap.Add(1)
	go func(waitSendRecieveMessageMap *sync.WaitGroup) {
		defer close(done)
		defer waitSendRecieveMessageMap.Done()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("close client!")
				return
			}
			sendMap[string(message)] = true
			if len(recvMap) >= NUM_REPEAT && len(sendMap) >= NUM_REPEAT {
				return
			}
		}
	}(waitSendRecieveMessageMap)

	// 송신부
	ticker := time.NewTicker(1)
	defer ticker.Stop()
	var i int
	for i < NUM_REPEAT {
		select {
		case <-done:
			return
		case tick := <-ticker.C:
			i++
			timeString := tick.String() + strconv.Itoa(i)
			recvMap[timeString] = true
			// timeBytes := []byte(timeString)
			// timeBytesSize := len(timeBytes)
			// data := make([]byte, 1024)
			// fmt.Println("send : ", len(data))
			err := c.WriteMessage(websocket.TextMessage, []byte(timeString))
			if err != nil {
				fmt.Println("write:", err)
				return
			}
		case <-interrupt:
			fmt.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(
					websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

	waitSendRecieveMessageMap.Wait()
	fmt.Printf("%d) recvMap %d \n", clientId, len(recvMap))
	fmt.Printf("%d) sendMap %d \n", clientId, len(sendMap))
	t.Logf("%d) recvMap %d \n", clientId, len(recvMap))

	if len(recvMap) != len(sendMap) {
		fmt.Printf("not equal %d %d \n", len(recvMap), len(sendMap))
		return
	}
	for k := range recvMap {
		if _, ok := sendMap[k]; !ok {
			fmt.Printf("not found %s \n", k)
			return
		}
	}
	fmt.Printf("success %d \n", clientId)
	t.Logf("success %d \n", clientId)
}

// var addr = flag.String("addr", "7a4fe220e12bb0c312a47e3885990a10.hi:7642", "http service address")
// func main() {
// 	res, _ := lib.EncryptSubdomain("0.0.0.0:8080")
// 	fmt.Println(res)
// 	res, _ = lib.EncryptSubdomain("0.0.0.0:8081")
// 	fmt.Println(res)
// 	res, _ = lib.EncryptSubdomain("0.0.0.0:8082")
// 	fmt.Println(res)
// }
