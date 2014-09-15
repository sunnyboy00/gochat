package main

// 聊天服务器(一般多个),客户端连接到此

import (
	"fmt"
	"net"
)

import (
	"share"
	"tcpserver"
	"tcpserver/endpoint"
	"types"
)

type Bot struct {
	endpoint.EndPoint

	User    *types.User
	Manager *TCPServerManager
}

func (bot *Bot) OnConnectionLost(err error) {
	fmt.Println("Connection Lost:", err.Error())

	bot.Ctrl <- false
	if bot.User.UID > 0 {
		//delete(bot.Manager.Clients, bot.User.UID)
		share.Clients.Delete(bot.User.UID)
	}
}

func (bot *Bot) Handle() {
	for {
		select {
		case data := <-bot.RecvBox:
			fmt.Println("Recv:", string(data))
			ack, err := HandleNetProto(bot, data)
			if err != nil {
				// 断开连接
				fmt.Println(err.Error())
				bot.Conn.Close()
				return
			}
			if ack != nil {
				bot.PutData(ack)
			}
		case data := <-bot.User.MQ: // internal IPC
			fmt.Println("MQ:", data)
		}
	}
}

type TCPServerManager struct {
	Address string
	//Clients map[uint32]*Bot // 这个应该加锁,如果是多核的话
}

func (m *TCPServerManager) connectionHandler(conn *net.TCPConn) {
	bot := &Bot{}
	bot.Init(conn, 180, 16, 12)
	bot.InitCBs(bot.OnConnectionLost, nil, nil)
	bot.Manager = m
	user := types.NewUser(8)
	bot.User = user

	//m.Clients[bot.User.UID] = bot
	share.Clients.Set(bot.User.UID, user)

	go bot.Handle()
	bot.Start()
}

func (m *TCPServerManager) Start() {
	server := tcpserver.NewStreamServer(m.Address, m.connectionHandler)
	server.Start()
}

func main() {
	manager := &TCPServerManager{Address: ":7005"}
	//manager.Clients = make(map[uint32]*Bot, 1000)
	go manager.Start()

	waiting := make(chan bool)
	<-waiting

}
