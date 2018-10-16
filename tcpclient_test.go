package gnet

import (
	"fmt"
	"testing"
)

type ClientConnListener struct{}

func (l *ClientConnListener) OnData(c *TcpConnection, buf []byte) {
	fmt.Println("client recv from ", c.RemoteAddr(), ":", string(buf))
}

func (l *ClientConnListener) OnClose(c *TcpConnection) {
	fmt.Println("client connection closed ", c.RemoteAddr())
}

func TestClient(t *testing.T) {
	_, err := NewServer()
	if err != nil {
		t.Failed()
	}

	conn, err := DialTcp("127.0.0.1:7001")

	if err != nil {
		t.Failed()
	}

	l := &ClientConnListener{}
	conn.RegisterListener(l)
	conn.StartHandle()
	conn.Write([]byte("test"))

	data := <-dataChan
	fmt.Println("recv from server data chan:", string(data))
	if string(data) != "test" {
		t.Failed()
	}

	conn.Close()
}
