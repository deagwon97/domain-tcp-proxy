package proxy

import (
	"io"
	"net"
)

type Tunnel struct {
	done     chan int
	appConn  net.Conn
	userConn net.Conn
	userBuf  io.Reader
}

func NewTunnel(
	appConn net.Conn,
	userConn net.Conn,
	userBuf io.Reader) *Tunnel {
	return &Tunnel{
		done:     make(chan int),
		appConn:  appConn,
		userConn: userConn,
		userBuf:  userBuf,
	}
}

func (tu *Tunnel) Stop() {
	tu.appConn.Close()
	tu.userConn.Close()
}

func (tu *Tunnel) copy(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
	tu.Stop()
}

func (tu *Tunnel) copyApp2User() {
	// appConn에서 읽어서 userBuf으로 쓰기
	// 에러 혹은 EOF 받기 전까지 blocking
	tu.copy(tu.appConn, tu.userBuf)
	// 사용자-중계서버 connection 끊어짐
	// channel에 1 push
	tu.done <- 1
}
func (tu *Tunnel) copyUser2App() {
	// userConn에서 읽어서 appConn으로 쓰기
	// 에러 혹은 EOF 받기 전까지 blocking
	tu.copy(tu.userConn, tu.appConn)
	// 중계서버-app connection 끊어짐
	// channel에 1 push
	tu.done <- 1
}

func (tu *Tunnel) Handle() {
	go tu.copyApp2User()
	go tu.copyUser2App()
	// tu.done 채널에 두 번 값이 push(혹은 send)하면 Handle() 함수 종료
	// 그 전까지 blocking
	<-tu.done
	<-tu.done
}
