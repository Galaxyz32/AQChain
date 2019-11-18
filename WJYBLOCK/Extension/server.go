package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	//"github.com/davecgh/go-spew/spew"
	"net"
	"strconv"
	"strings"
	"time"
)

var ServerIsNode *bool

var counter = 1
var sendaddr string
var IDPort = "19099"
var AddrPort = "19093"
var IDMutex = &sync.Mutex{}

type NodeStruct struct {
	ID       int
	Addr     string
	IsOnline bool
	PubKey   []byte
}

var NodeList []NodeStruct

func AppendNodes(id string, addr string, Key []byte) {
	num, _ := strconv.Atoi(id)
	newNode := NodeStruct{num, addr, true, Key}
	NodeList = append(NodeList, newNode)
}

func AddrBookListen() {

	AddrBookListener, _ := net.Listen("tcp", ":"+AddrPort)
	for {
		AddrBook, _ := AddrBookListener.Accept()
		go AddrBookHandler(AddrBook)
	}

}
func AddrBookHandler(conn net.Conn) {

	addr := conn.RemoteAddr().String()
	ip := strings.Split(addr, ":")
	//PublicIP := ip[0]
	//fmt.Println(PublicIP)

	if ip[0] == "127.0.0.1" && *ServerIsNode == true {
		ip[0] = GetExternalIP()
		ip[0] = strings.Replace(ip[0], "\n", "", -1)
	}

	NLbytes, _ := json.Marshal(NodeList)
	fmt.Fprintln(conn, string(NLbytes))
	sendaddr, _ = bufio.NewReader(conn).ReadString(24)
	if sendaddr == "" {
		conn.Close()
		return
	}

	AddrKey := strings.Split(sendaddr, "|")
	//spew.Dump(AddrKey)

	str := strings.Split(AddrKey[1], "/")
	AddrKey[1] = strings.Replace(AddrKey[1], str[2], ip[0], -1)

	fmt.Println("新上线的节点序号是" + AddrKey[0])
	fmt.Println("接收到：端口" + str[4] + "，地址" + str[2])

	AddrKey[2] = strings.Trim(AddrKey[2], string(24))

	fmt.Println(AddrKey[1])
	AppendNodes(AddrKey[0], AddrKey[1], []byte(AddrKey[2]))
	fmt.Print("已上线，ID是")
	fmt.Println(NodeList[len(NodeList)-1].ID)

	conn.Close()

}

func main() {

	ServerIsNode = flag.Bool("s", false, "ID服务器上是否也在运行节点")
	flag.Parse()

	//初始化用户表
	node := NodeStruct{0, "", false, nil}
	NodeList = append(NodeList, node)
	go idsvr()
	go AddrBookListen()
	go isAliveTest()
	select {}

}

func idsvr() {
	listener, _ := net.Listen("tcp", ":"+IDPort)
	for {
		conn, _ := listener.Accept()
		//每个客户端一个goroutine
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	IDMutex.Lock()
	fmt.Fprintln(conn, counter)
	fmt.Print("序号")
	fmt.Print(counter)
	fmt.Println("已分配")
	counter++
	conn.Close()
	IDMutex.Unlock()
}

func GetExternalIP() string {
	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	//s := buf.String()
	return string(content)
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
				if err != nil {
					NodeList[i].IsOnline = false
					log.Print("节点 " + str[4] + "已离线")
				} else {
					conn.Close()
				}
			}

		}
	}
}
