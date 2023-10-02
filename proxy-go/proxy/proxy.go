package proxy

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	done chan int
	host string
}

func NewServer(host string) *Server {
	return &Server{
		done: make(chan int),
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
	return nil
}

func (s *Server) Stop() {
	<-s.done
	close(s.done)
}

func (s *Server) run(listener *net.TCPListener) {
	for {
		select {
		case <-s.done:
			return
		default:
			// 새로운 user의 요청을 기다림
			userConn, err := listener.AcceptTCP()
			if err != nil {
				return
			}

			// 새로운 요청이 들어오고 나서의 과정은
			// go routine을 생성해서 수행
			go func(userConn *net.TCPConn) {
				// http 요청을 파싱해서 목적지 AppHost를 구함
				userReqeustReader := bufio.NewReader(userConn)
				userReqeust, err := http.ReadRequest(userReqeustReader)
				if err != nil {
					return
				}
				appHost := ParseAppHost(userReqeust)

				// 사용자 http 요청을 파싱하면서 원래 데이터가 사라짐
				// 원래 사용자의 http 요청을 복원
				userRequestBuf := bytes.NewBuffer([]byte{})
				userRequestBytes, _ := httputil.DumpRequest(userReqeust, true)
				userRequestBuf.Write(userRequestBytes)
				// http reqeust data io.Reader 와 userConn io.Reader를 결합
				userBuf := io.MultiReader(userRequestBuf, userConn)
				// 목적지 appHost와 tcp connection 연결
				appConn, _ := net.Dial("tcp", appHost)
				// tunnel 객체 생성
				tu := NewTunnel(
					appConn,
					userConn,
					userBuf,
				)
				// <사용자-중계서버> ---- <중계서버-목적지앱>
				// 사이를 tcp 터널링
				tu.Handle()
			}(userConn)

		}
	}
}
