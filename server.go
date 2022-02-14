package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnLineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//监听Message广播消息channel的goroutine，一旦有消息就发送给全部现在User
func (pthis *Server) ListenMessage() {
	for {
		msg := <-pthis.Message
		//将msg发送给全部在线用户
		pthis.mapLock.Lock()
		for _, cli := range pthis.OnLineMap {
			cli.C <- msg
		}
		pthis.mapLock.Unlock()
	}
}

func (pthis *Server) Handler(conn net.Conn) {
	user := NewUser(conn)

	//用户上线，将用户加入到OnLineMap中
	pthis.mapLock.Lock()
	pthis.OnLineMap[user.Name] = user
	pthis.mapLock.Unlock()

	//广播当前用户上线消息
	pthis.BroadCast(user, "已上线")

	//当前handle阻塞
	select {}
}

//广播消息的方法
func (pthis *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	pthis.Message <- sendMsg
}

//启动服务器的接口
func (pthis *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", pthis.Ip, pthis.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	// close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go pthis.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go pthis.Handler(conn)

	}

	//accept

	//do handler

	//close listen socket

}
