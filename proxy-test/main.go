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

func main() {
	NUM_OF_SERVER, _ := strconv.Atoi(os.Args[1])
	NUM_OF_CLIENT, _ := strconv.Atoi(os.Args[2])
	NUM_REPEAT, _ := strconv.Atoi(os.Args[3])
	DATA_SIZE, _ := strconv.Atoi(os.Args[4])
	run(NUM_OF_SERVER, NUM_OF_CLIENT, NUM_REPEAT, DATA_SIZE)
}
