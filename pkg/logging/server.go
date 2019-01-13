package logging

import (
	"net"
	"fmt"
)

// Client holds the connection handle
type LogServer struct {
	conn net.Conn
}

func StartLogServer(port int) *LogServer {

	endpoint := fmt.Sprintf(":%d", port)

	ln, err := net.Listen("tcp", endpoint)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Please open another terminal window and execute the command 'kubeprov logs'")
  	conn, err := ln.Accept() //blocking
	if err != nil {
		fmt.Println(err)
	}

	return &LogServer {
		conn: conn,
	}
}

func (s *LogServer) SendLogMessage(msg string){
	s.conn.Write([]byte(msg + "\n"))
}