package gnet

import (
	"fmt"
	"testing"
)

func HandleNewConn(s *TcpServer) {
	for {
		c := <-s.NewConnChan()
		fmt.Println("new connection from ", c.RemoteAddr())
		l := &ServerConnListener{}
		c.RegisterListener(l)
		c.StartHandle()
	}
}

var dataChan = make(chan []byte)

type ServerConnListener struct{}

func (l *ServerConnListener) OnData(c *TcpConnection, buf []byte) {
	fmt.Println("server recv from ", c.RemoteAddr(), ":", string(buf))
	dataChan <- buf
	c.Write(buf)
}

func (l *ServerConnListener) OnClose(c *TcpConnection) {
	fmt.Println("server connection closed ", c.RemoteAddr())
}

func NewServer() (*TcpServer, error) {
	s := TcpServer{}
	err := s.Start(":7001")
	if err != nil {
		return nil, err
	}

	go HandleNewConn(&s)

	return &s, nil
}

func TestServer(t *testing.T) {
	_, err := NewServer()
	if err != nil {
		t.Failed()
	}

	data := <-dataChan
	fmt.Println("recv from server data chan:", string(data))
}
