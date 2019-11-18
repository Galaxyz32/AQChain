package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var counter = 1
var sendaddr string

type NodeStruct struct {
	ID       int
	Addr     string
	IsOnline bool
}

var NodeList []NodeStruct

func AppendNodes(addr string) {
	newNode := NodeStruct{counter - 1, addr, true}
	NodeList = append(NodeList, newNode)
}

func AddrBookListen() {

	AddrBookListener, _ := net.Listen("tcp", ":19093")
	for {
		AddrBook, _ := AddrBookListener.Accept()
		go AddrBookHandler(AddrBook)
	}

}

func AddrBookHandler(conn net.Conn) {

	bytes, _ := json.Marshal(NodeList)
	fmt.Fprintln(conn, string(bytes))
	sendaddr, _ = bufio.NewReader(conn).ReadString('\n')
	str := strings.Split(sendaddr, "/")
	fmt.Println("接收到：端口" + str[4] + "，地址" + str[2])

	sendaddr = strings.Replace(sendaddr, "\n", "", -1)
	AppendNodes(sendaddr)
	fmt.Print("已上线，ID是")
	fmt.Println(NodeList[len(NodeList)-1].ID)
	conn.Close()

}

func main() {

	//初始化用户表
	node := NodeStruct{0, "", false}
	NodeList = append(NodeList, node)
	go idsvr()
	go AddrBookListen()
	go isAliveTest()
	select {}

}

func idsvr() {
	listener, _ := net.Listen("tcp", ":19099")
	for {
		conn, _ := listener.Accept()
		//每个客户端一个goroutine
		go handleConn(conn)
		fmt.Print("序号")
		fmt.Print(counter)
		fmt.Println("已分配")
	}
}

func handleConn(conn net.Conn) {
	fmt.Fprintln(conn, counter)
	conn.Close()
	counter++
}

func isAliveTest() {
	fmt.Println("开始探测节点是否在线")
	for {
		time.Sleep(3e9)
		for i := 1; i < len(NodeList); i++ {
			if NodeList[i].IsOnline == true {
				str := strings.Split(NodeList[i].Addr, "/")
				port, _ := strconv.Atoi(str[4])
				addr := &net.TCPAddr{IP: net.ParseIP(str[2]), Port: port}
				conn, err := net.DialTCP("tcp", nil, addr)
				//fmt.Print("正在扫描" + str[4])
				if err != nil {
					NodeList[i].IsOnline = false
					fmt.Println(str[4] + "已经离线")
				} else {
					conn.Close()
				}
			}

		}
	}
}
