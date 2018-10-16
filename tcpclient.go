package gnet

import (
	"net"
)

func DialTcp(address string) (*TcpConnection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	tcpConn := NewTcpConnection(nil, conn)
	return tcpConn, nil
}
