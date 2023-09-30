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
	proxyHost := "0.0.0.0:9980"
	fmt.Printf("proxy server start on %s ...\n", proxyHost)
	s := proxy.NewServer(proxyHost)
	s.Start()
	defer s.Stop()
	select {}
}

func main() {
	// go recoverer(api.EncryptSubdomainApi)
	go recoverer(runProxy)
	select {}
}
