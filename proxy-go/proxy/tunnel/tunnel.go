package tunnel

import (
	"io"
	"net"
	"sync"
)

type Tunnel struct {
	done     chan struct{}
	appConn  net.Conn
	userConn net.Conn
	userBuf  io.Reader
}

func NewTunnel(
	appConn net.Conn,
	userConn net.Conn,
	userBuf io.Reader) *Tunnel {
	// appHost := appConn.RemoteAddr().String()
	// userHost := userConn.RemoteAddr().String()
	// fmt.Println("create tunnel: ", userHost, "<->", appHost)
	return &Tunnel{
		done:     make(chan struct{}),
		appConn:  appConn,
		userConn: userConn,
		userBuf:  userBuf,
	}
}

func (tu *Tunnel) Stop() {
	if tu.done == nil {
		return
	}
	// appHost := tu.appConn.RemoteAddr().String()
	// userHost := tu.userConn.RemoteAddr().String()
	// fmt.Println("close tunnel: ", userHost, "<->", appHost)
	close(tu.done)
	tu.done = nil
	if tu.appConn != nil {
		tu.appConn.Close()
		tu.appConn = nil
	}
	if tu.userConn != nil {
		tu.userConn.Close()
		tu.userConn = nil
	}
}

func (tu *Tunnel) Handle() {
	defer tu.Stop()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go tu.copyApp2User(wg)
	go tu.copyUser2App(wg)
	wg.Wait()
}

func (tu *Tunnel) copyApp2User(wg *sync.WaitGroup) {
	copy(tu.appConn, tu.userBuf, wg, tu)
}
func (tu *Tunnel) copyUser2App(wg *sync.WaitGroup) {
	copy(tu.userConn, tu.appConn, wg, tu)
}

func copy(dst io.Writer, src io.Reader,
	wg *sync.WaitGroup, tu *Tunnel) {
	defer wg.Done()
	defer tu.Stop()
	select {
	case <-tu.done:
		return
	default:
		// minimum packet size of TCP is 20 bytes
		// buffer := make([]byte, 20)
		if _, err := io.Copy(dst, src); err != nil {
			tu.Stop()
			return
		}
	}
}
