package tunnel

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

func (tu *Tunnel) Handle() {
	go tu.copyApp2User()
	go tu.copyUser2App()
	<-tu.done
	<-tu.done
}

func (tu *Tunnel) copyApp2User() {
	tu.copy(tu.appConn, tu.userBuf)
	tu.done <- 1
}
func (tu *Tunnel) copyUser2App() {
	tu.copy(tu.userConn, tu.appConn)
	tu.done <- 1
}

func (tu *Tunnel) copy(dst io.Writer, src io.Reader) {
	io.Copy(dst, src)
	tu.Stop()
}
