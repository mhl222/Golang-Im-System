package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// NewServer 创建server，工厂方法
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// ListenMessage 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
func (receiver *Server) ListenMessage() {
	fmt.Println("消息广播协程启动")
	for {
		msg := <-receiver.Message
		receiver.mapLock.Lock()
		for _, cli := range receiver.OnlineMap {
			cli.C <- msg
		}
		receiver.mapLock.Unlock()
	}
}

// BroadCast 用户：广播消息
func (receiver *Server) BroadCast(user *User, msg string) {
	//sendMsg := fmt.Sprintf("[%s] %s: %s\n", user.Addr, user.Name, msg)
	sendMsg := "[" + user.Addr + "]" + " " + user.Name + ": " + msg + "\r\n"
	receiver.Message <- sendMsg
}

func (receiver *Server) handler(con net.Conn) {
	fmt.Println("链接建立成功")

	user := NewUser(con, receiver)

	user.Online()

	// 接受客户端消息
	go func(conn net.Conn) {
		buf := make([]byte, 4096)
		for {
			n, err := con.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("消息接收失败", err)
				return
			}

			// 提取用户消息
			msg := string(buf[:n-1])
			// 用户消息广播
			user.DoMessage(msg)
		}
	}(con)
}

func (receiver *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", receiver.Ip, receiver.Port))
	if err != nil {
		fmt.Println("net.Listener err", err)
		return
	}
	fmt.Printf("server run at %s:%d\n", receiver.Ip, receiver.Port)

	// close listen socket
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("net.Listener close err", err)
		}
	}(listener)

	// 启动监听Message广播消息channel的goroutine
	go receiver.ListenMessage()
	for {

		//accept
		conn, err := listener.Accept()
		if err != nil {

			fmt.Println("listener accept err", err)
			continue
		}
		//do handler
		go receiver.handler(conn)

	}

}
