package main

import (
	"fmt"
	"proxy-go/proxy"
)

func recoverer(f func()) {
	defer func() {
		if err := recover(); err != nil {
			recoverer(f)
		}
	}()
	f()
}

func runProxy() {
	fmt.Println("run proxy server")
	proxyHost := "0.0.0.0:8080"
	s := proxy.NewServer(proxyHost)
	s.Start()
	defer s.Stop()
	select {}
}

func main() {
	recoverer(runProxy)
	select {}
}
