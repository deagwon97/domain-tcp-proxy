package proxy

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"proxy-go/lib"
	"proxy-go/proxy/tunnel"
)

type Server struct {
	done chan struct{}
	host string
}

func NewServer(host string) *Server {
	return &Server{
		done: make(chan struct{}),
		host: host,
	}
}

func (s *Server) Start() error {
	listenAddr, err := net.ResolveTCPAddr("tcp", s.host)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return err
	}
	go s.run(listener)
	// panic("test")
	return nil
}

func (s *Server) Stop() {
	close(s.done)
}

func (s *Server) run(listener *net.TCPListener) {
	for {
		select {
		case <-s.done:
			return
		default:
			userConn, err := listener.AcceptTCP()
			if err != nil {
				return
			}
			userReqeustReader := bufio.NewReader(userConn)
			// userReqeustReader.Seek(0, 0)

			// parse http reqeust
			userReqeust, err := http.ReadRequest(userReqeustReader)
			if err != nil {
				return
			}
			appHost := lib.ParseAppHost(userReqeust)

			// recover http reqeust bytes data to io.Reader
			userRequestBytes, _ := httputil.DumpRequest(userReqeust, true)
			userRequestBuf := bytes.NewBuffer([]byte{})
			userRequestBuf.Write(userRequestBytes)

			// combine http reqeust data io.Reader and userConn io.Reader
			userBuf := io.MultiReader(userRequestBuf, userConn)
			appConn, _ := net.Dial("tcp", appHost)
			// proxy 객체 생성
			tu := tunnel.NewTunnel(
				appConn,
				userConn,
				userBuf,
			)
			tu.Handle()
			// if err == nil {
			// 	go tu.Handle()
			// } else {
			// 	fmt.Println("Error accepting conn")
			// }
		}
	}
}
