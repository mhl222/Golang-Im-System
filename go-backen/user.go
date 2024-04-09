package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser 创建工厂
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		userAddr,
		userAddr,
		make(chan string),
		conn,
		server,
	}
	go user.ListenMessage()
	return user
}

// ListenMessage 监听当前user channel的消息，并发送给客户端
func (receiver User) ListenMessage() {
	for {
		msg := <-receiver.C
		_, err := receiver.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("conn write error:", err)
		}
	}
}

// Online 用户上线
func (receiver *User) Online() {
	// 用户上线，将用户加入OnlineMap
	receiver.server.mapLock.Lock()
	receiver.server.OnlineMap[receiver.Name] = receiver
	receiver.server.mapLock.Unlock()
	// 广播用户上线消息
	receiver.server.BroadCast(receiver, "Get online\r\n")
}

// Offline 用户下线
func (receiver *User) Offline() {

	// 将用户从OnlineMap删除
	receiver.server.mapLock.Lock()
	delete(receiver.server.OnlineMap, receiver.Name)
	receiver.server.mapLock.Unlock()

	// 广播用户下线消息
	receiver.server.BroadCast(receiver, "off  online\r\n")

}

// DoMessage 用户处理消息
func (receiver *User) DoMessage(msg string) {
	receiver.server.BroadCast(receiver, msg)

}
