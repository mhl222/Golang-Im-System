package main

import (
	"fmt"
	"net"
	"strings"
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
func (receiver *User) ListenMessage() {
	for {
		msg := <-receiver.C
		_, err := receiver.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("conn write error:", err)
			return
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

func (receiver *User) SendMsg(msg string) {
	_, err := receiver.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn write error:", err)
		return
	}
}

// DoMessage 用户处理消息
func (receiver *User) DoMessage(msg string) {

	if msg == "who" {
		// 查询当前在线用户
		receiver.server.mapLock.Lock()
		for _, user := range receiver.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "Online\r\n"
			receiver.SendMsg(onlineMsg)
		}
		receiver.server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[0:7] == "rename|" {
		// 修改用户名
		// 消息格式 rename|张三
		newName := strings.Split(msg, "|")[1]
		// 判断name是否存在
		_, ok := receiver.server.OnlineMap[newName]
		if ok {
			receiver.SendMsg("name already exist\r\n")
		} else {

			receiver.server.mapLock.Lock()
			delete(receiver.server.OnlineMap, receiver.Name)
			receiver.server.OnlineMap[newName] = receiver
			receiver.server.mapLock.Unlock()
			receiver.Name = newName
			receiver.SendMsg("rename success\r\n")
		}
	} else {
		receiver.server.BroadCast(receiver, msg)
	}

}
