package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser 创建工厂
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		userAddr,
		userAddr,
		make(chan string),
		conn,
	}
	go user.ListenMessage()
	return user
}
func (receiver User) ListenMessage() {
	for {
		msg := <-receiver.C
		_, err := receiver.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("conn write error:", err)
		}
	}
}
