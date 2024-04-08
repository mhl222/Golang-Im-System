package main

import (
	"fmt"
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

// 创建server，工厂方法
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听Message广播消息channel的goroutine，一旦有消息就发送给全部的在线User
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

// 广播消息
func (receiver *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s] %s: %s\n", user.Addr, user.Name, msg)
	receiver.Message <- sendMsg
}

func (receiver *Server) handler(con net.Conn) {
	fmt.Println("链接建立成功")

	// 用户上线，将用户加入OnlineMap
	user := NewUser(con)
	receiver.mapLock.Lock()
	receiver.OnlineMap[user.Name] = user
	receiver.mapLock.Unlock()
	// 广播用户上线消息
	receiver.BroadCast(user, "已上线")

}

func (receiver *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%receiver:%d", receiver.Ip, receiver.Port))
	if err != nil {
		fmt.Println("net.Listener err", err)
		return
	}
	fmt.Printf("server run at %receiver:%d\n", receiver.Ip, receiver.Port)

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
