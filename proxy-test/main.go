package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"proxy-test/lib"
	"proxy-test/server"
	"strconv"
	"sync"
	"time"
<<<<<<< HEAD
=======

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/blowfish"
>>>>>>> e4ea4cf (update)
)

const (
	APP_BASE_PORT       = 13001
	TUNNEL_HOST_POSTFIX = ".service.com"
	MID_SERVER_PORT     = 9980
)

type Response struct {
	EncryptHost string `json:"encrypted_host"`
}

func run(
	NUM_OF_SERVER int,
	NUM_OF_CLIENT int,
	NUM_REPEAT int,
	DATA_SIZE int,
) {
	waitWebSocketServerReady := &sync.WaitGroup{}

<<<<<<< HEAD
	//	run servers
	waitWebSocketServerReady.Add(NUM_OF_SERVER)
	for i := 0; i < NUM_OF_SERVER; i++ {
		ip := "0.0.0.0"
		port := strconv.Itoa(APP_BASE_PORT + i)
		appHost := ip + ":" + port
		go server.RunServer(waitWebSocketServerReady, appHost)
		host := fmt.Sprintf("%s:%s", ip, port)
		encryptSubdomain, _ := lib.EncryptSubdomain(host)
		domain := encryptSubdomain + TUNNEL_HOST_POSTFIX
		lib.AddHostEntry(domain, "0.0.0.0")
	}
	// wait for all servers ready
	waitWebSocketServerReady.Wait()

	waitAllClientsEnd := &sync.WaitGroup{}
	waitAllClientsEnd.Add(NUM_OF_CLIENT * NUM_OF_SERVER)
	startAt := time.Now()
	for nc := 0; nc < NUM_OF_CLIENT; nc++ {
		for ns := 0; ns < NUM_OF_SERVER; ns++ {
			appHost := fmt.Sprintf("0.0.0.0:%d", APP_BASE_PORT+ns)
			go server.RunClient(waitAllClientsEnd,
				ns+1,
				MID_SERVER_PORT,
				appHost,
				TUNNEL_HOST_POSTFIX,
				NUM_REPEAT,
				DATA_SIZE,
			)
		}
	}

	waitAllClientsEnd.Wait()
	endAt := time.Now()
	duration := endAt.Sub(startAt).Seconds()
	fmt.Printf("%d, %d, %d, %d, %f \n", NUM_OF_SERVER, NUM_OF_CLIENT, NUM_REPEAT, DATA_SIZE, duration)
}

=======
>>>>>>> e4ea4cf (update)
func main() {
	NUM_OF_SERVER, _ := strconv.Atoi(os.Args[1])
	NUM_OF_CLIENT, _ := strconv.Atoi(os.Args[2])
	NUM_REPEAT, _ := strconv.Atoi(os.Args[3])
	DATA_SIZE, _ := strconv.Atoi(os.Args[4])
	run(NUM_OF_SERVER, NUM_OF_CLIENT, NUM_REPEAT, DATA_SIZE)
}
<<<<<<< HEAD
=======

func updateEtcHosts() {
	f, err := os.OpenFile("/etc/hosts",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString("\n7a4fe220e12bb0c312a47e3885990a10.service.com:0.0.0.0"); err != nil {
		log.Println(err)
	}
	if _, err := f.WriteString("\n7a4fe220e12bb0c3ba67abb5cef3a8c0.service.com:0.0.0.0"); err != nil {
		log.Println(err)
	}
	if _, err := f.WriteString("\n7a4fe220e12bb0c344cb29f16bbff67a.service.com:0.0.0.0"); err != nil {
		log.Println(err)
	}

}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func EncryptSubdomain(host string) (subDomain string, err error) {
	key := []byte("thisissecretkey")
	block, err := blowfish.NewCipher(key)
	if err != nil {
		log.Println(err)
	}
	mode := ecb.NewECBEncrypter(block)
	plaintext := []byte(host)
	plaintext = PKCS5Padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	subDomain = hex.EncodeToString(ciphertext)
	return subDomain, err
}

// test http based proxy server performance
func TestWebScoket() {
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
	waitAllClientsEnd.Add(NUM_OF_CLIENT * NUM_OF_SERVER)
	for i := 0; i < NUM_OF_CLIENT; i++ {
		go runClient(waitAllClientsEnd, 1, appHost1)
		go runClient(waitAllClientsEnd, 2, appHost2)
		go runClient(waitAllClientsEnd, 3, appHost3)
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

func runClient(waitClientEnd *sync.WaitGroup, clientId int, appHost string) {
	defer waitClientEnd.Done()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	subdomain, _ := EncryptSubdomain(appHost)
	tunnelHost := subdomain + TUNNEL_HOST_POSTFIX
	fmt.Println("Connecting to " + tunnelHost)
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

}

// var addr = flag.String("addr", "7a4fe220e12bb0c312a47e3885990a10.service.com:7642", "http service address")
// func main() {
// 	res, _ := lib.EncryptSubdomain("0.0.0.0:8080")
// 	fmt.Println(res)
// 	res, _ = lib.EncryptSubdomain("0.0.0.0:8081")
// 	fmt.Println(res)
// 	res, _ = lib.EncryptSubdomain("0.0.0.0:8082")
// 	fmt.Println(res)
// }
>>>>>>> e4ea4cf (update)
