package gnet

import (
	"math"
	"net"
)

type TcpServer struct {
	connNum     int
	connNumChan chan int
	listener    net.Listener
	newConnChan chan *TcpConnection
	isRunning   bool
}

func (s *TcpServer) Start(addr string) error {
	if s.isRunning {
		return ErrServerAlreadyStart
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = listener
	s.connNumChan = make(chan int)
	s.newConnChan = make(chan *TcpConnection)
	go s.calcConnNum()
	go s.acceptClients()

	s.isRunning = true

	return nil
}

func (s *TcpServer) ConnNum() int {
	return s.connNum
}

func (s *TcpServer) NewConnChan() chan *TcpConnection {
	return s.newConnChan
}

func (s *TcpServer) calcConnNum() {
	for {
		i := <-s.connNumChan
		if i == math.MaxInt32 {
			break
		}
		s.connNum += i
	}
}

func (s *TcpServer) handleConnection(conn net.Conn) {
	newConn := NewTcpConnection(s, conn)
	s.connNumChan <- 1
	s.newConnChan <- newConn
}

func (s *TcpServer) acceptClients() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		s.handleConnection(conn)
	}
}
