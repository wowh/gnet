package gnet

import (
	"net"
)

type TcpConnection struct {
	server       *TcpServer
	conn         net.Conn
	connListener ConnListener
	isShutdown   bool
	writeChan    chan []byte
	closeChan    chan struct{}
}

type ConnListener interface {
	OnData(*TcpConnection, []byte)
	OnClose(*TcpConnection)
}

func NewTcpConnection(s *TcpServer, c net.Conn) *TcpConnection {
	conn := &TcpConnection{}
	conn.server = s
	conn.conn = c
	conn.isShutdown = false
	conn.writeChan = make(chan []byte)
	conn.closeChan = make(chan struct{})
	return conn
}

func (c *TcpConnection) RegisterListener(listener ConnListener) {
	c.connListener = listener
}

func (c *TcpConnection) StartHandle() error {
	if c.connListener == nil {
		return ErrConnListenerNotRegister
	}

	go c.connRead()
	go c.connWrite()
	go c.handleConnClose()
	return nil
}

func (c *TcpConnection) Write(buf []byte) error {
	for {
		select {
		case <-c.closeChan:
			return ErrTcpConnectionAlreadyClosed
		case c.writeChan <- buf:
			return nil
		}
	}
}

func (c *TcpConnection) Close() error {
	if !c.isShutdown {
		return c.conn.Close()
	}

	return ErrTcpConnectionAlreadyClosed
}

func (c *TcpConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TcpConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *TcpConnection) connRead() {
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf)
		if err != nil {
			c.conn.Close()
			break
		}

		c.connListener.OnData(c, buf[:n])
	}

	close(c.closeChan)
}

func (c *TcpConnection) connWrite() {
	for {
		select {
		case buf := <-c.writeChan:
			_, err := c.conn.Write(buf)
			if err != nil {
				// close connnection to notify connRead
				c.conn.Close()
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *TcpConnection) handleConnClose() {
	<-c.closeChan
	if !c.isShutdown {
		c.isShutdown = true
		c.connListener.OnClose(c)
		if c.server != nil {
			c.server.connNumChan <- -1
		}
	}
}
