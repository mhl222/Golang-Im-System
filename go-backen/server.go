package main

import (
	"fmt"
	"net"
)

// 创建server，工厂方法
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

type Server struct {
	Ip   string
	Port int
}

func (s *Server) handler(con net.Conn) {
	fmt.Println("链接建立成功")

}

func (s *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listener err", err)
		return
	}
	fmt.Printf("server run at %s:%d\n", s.Ip, s.Port)

	// close listen socket
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("net.Listener close err", err)
		}
	}(listener)

	for {

		//accept
		conn, err := listener.Accept()
		if err != nil {

			fmt.Println("listener accept err", err)
			continue
		}
		//do handler
		go s.handler(conn)

	}

}
