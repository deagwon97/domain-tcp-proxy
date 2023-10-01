package main

import (
	"fmt"
	"os"
	"proxy-test/lib"
	"proxy-test/server"
	"strconv"
	"sync"
	"time"
)

const (
	APP_BASE_PORT       = 7080
	TUNNEL_HOST_POSTFIX = ".service.com"
)

type Response struct {
	EncryptHost string `json:"encrypted_host"`
}

func run(
	NUM_OF_SERVER int,
	NUM_OF_CLIENT int,
	NUM_REPEAT int,
	PACKET_SIZE int,
) {
	waitWebSocketServerReady := &sync.WaitGroup{}

	// fmt.Println("start test")

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

	//time start
	startAt := time.Now()
	// run clients
	waitAllClientsEnd := &sync.WaitGroup{}
	waitAllClientsEnd.Add(NUM_OF_CLIENT * NUM_OF_SERVER)
	for nc := 0; nc < NUM_OF_CLIENT; nc++ {
		for ns := 0; ns < NUM_OF_SERVER; ns++ {
			appHost := fmt.Sprintf("0.0.0.0:%d", APP_BASE_PORT+ns)
			// fmt.Println("run client")
			go server.RunClient(waitAllClientsEnd, ns+1, appHost, TUNNEL_HOST_POSTFIX, NUM_REPEAT, PACKET_SIZE)
		}
	}
	// wait for all test end
	waitAllClientsEnd.Wait()
	// time end
	endAt := time.Now()
	delta := endAt.Sub(startAt)
	fmt.Printf("%d, %d, %d, %d, %f \n", NUM_OF_SERVER, NUM_OF_CLIENT, NUM_REPEAT, PACKET_SIZE, delta.Seconds())

}

func main() {
	NUM_OF_SERVER, _ := strconv.Atoi(os.Args[1])
	NUM_OF_CLIENT, _ := strconv.Atoi(os.Args[2])
	NUM_REPEAT, _ := strconv.Atoi(os.Args[3])
	PACKET_SIZE, _ := strconv.Atoi(os.Args[4])
	run(NUM_OF_SERVER, NUM_OF_CLIENT, NUM_REPEAT, PACKET_SIZE)
}
