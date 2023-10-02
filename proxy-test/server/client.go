package server

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"proxy-test/lib"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func RunClient(waitClientEnd *sync.WaitGroup,
	clientId int,
	midServerPort int,
	appHost string,
	TUNNEL_HOST_POSTFIX string,
	NUM_REPEAT int,
	DATA_SIZE int) {
	defer waitClientEnd.Done()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	subdomain, _ := lib.EncryptSubdomain(appHost)
	tunnelHost := subdomain + TUNNEL_HOST_POSTFIX + ":" + strconv.Itoa(midServerPort)
	u := url.URL{Scheme: "ws", Host: tunnelHost, Path: "/echo"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	sendCount := 0
	recvCount := 0

	// 수신부
	done := make(chan struct{})
	waitSendRecieveMessageMap := &sync.WaitGroup{}
	waitSendRecieveMessageMap.Add(1)
	go func(waitSendRecieveMessageMap *sync.WaitGroup) {
		defer close(done)
		defer waitSendRecieveMessageMap.Done()
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
			recvCount++
			if recvCount >= NUM_REPEAT {
				return
			}
		}
	}(waitSendRecieveMessageMap)

	// 송신부
	ticker := time.NewTicker(10)
	defer ticker.Stop()
	var i int
	for i < NUM_REPEAT {
		select {
		case <-done: // 종료
			return
		case <-interrupt: // 종료 메시지가 들어올 경우
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
		case <-ticker.C: // 1 nano second 단위로 전송
			i++
			data := make([]byte, DATA_SIZE)
			err := c.WriteMessage(websocket.BinaryMessage, data)
			sendCount++
			if err != nil {
				fmt.Println("write:", err)
				return
			}

		}
	}
	waitSendRecieveMessageMap.Wait()
}
