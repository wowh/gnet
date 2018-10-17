# gnet
a go network library

# usage
## start server
```go
func NewServer() (*TcpServer, error) {
	s := TcpServer{}
	err := s.Start(":7001")
	if err != nil {
		return nil, err
	}

	return &s, nil
}
```
## handle connection

### implement a struct with follow OnData and OnClose function
```go
type ServerConnListener struct{}

func (l *ServerConnListener) OnData(c *TcpConnection, buf []byte) {
	fmt.Println("server recv from ", c.RemoteAddr(), ":", string(buf))
	dataChan <- buf
	c.Write(buf)
}

func (l *ServerConnListener) OnClose(c *TcpConnection) {
	fmt.Println("server connection closed ", c.RemoteAddr())
}
```

### implement a function to handle new connection from server
```go
func HandleNewConn(s *TcpServer) {
	for {
		c := <-s.NewConnChan()
		fmt.Println("new connection from ", c.RemoteAddr())
		l := &ServerConnListener{}
		c.RegisterListener(l)
		c.StartHandle()
	}
}
```

### 
