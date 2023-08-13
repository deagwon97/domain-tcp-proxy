package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/blowfish"
)

const (
	NUM_OF_SERVER = 3 // fixed
	NUM_REPEAT    = 100
	NUM_OF_CLIENT = 100
	// PACKET_SIZE         = 1024 // minimum 128
	TUNNEL_HOST_POSTFIX = ".service.com:8080"
)

var upgrader = websocket.Upgrader{} // use default options

func main() {
	TestWebScoket()
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

func hostExists(host string) (bool, error) {
	file, err := os.Open("/etc/hosts")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, host) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func addHostEntry(host, ip, port string) error {
	entry := fmt.Sprintf("%s\t%s:%s", ip, host, port)

	// Check if the entry already exists in the hosts file
	hostsFile, err := os.ReadFile("/etc/hosts")
	if err != nil {
		return err
	}

	if strings.Contains(string(hostsFile), entry) {
		return fmt.Errorf("Entry already exists in /etc/hosts")
	}

	// Append the new entry to the hosts file
	file, err := os.OpenFile("/etc/hosts", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, entry)
	if err != nil {
		return err
	}

	return nil
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

	sub1 := EncryptSubdomain("appHost1")

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
