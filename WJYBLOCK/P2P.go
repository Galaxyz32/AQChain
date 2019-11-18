package main

import (
	"bufio"
	"bytes"
	"container/list"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"log"
	csnet "net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//用于P2P连接的全局变量
var distAddr string
var ports int
var myaddr string
var ServerIP *string
var send = make(chan string, 100)
var recv = make(chan string, 100)

var SendList *list.List

var P2PAddrBookPort = "19093"
var P2PIDPort = "19099"

// 地址本
type NodeStruct struct {
	ID       int
	Addr     string
	IsOnline bool
	PubKey   []byte
}

var NodeList []NodeStruct

//读写锁
var p2pmutex = &sync.Mutex{}

//指定监听
func makeBasicHost(listenPort int) host.Host {

		opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		//libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/10.18.19.211/tcp/%d", listenPort)),
	}
		basicHost, _ := libp2p.New(context.Background(), opts...)

	// 生成对外服务地址,端口+1
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	//basicHost.Addrs()是当前监听的地址列表

	for _, addr := range basicHost.Addrs() {
		if strings.SplitN(addr.String(), "/", 3)[1] != "p2p-circuit" { //地址分三段，第一段不是中继标志
			//fullAddr是协议封装后的地址
			fullAddr := addr.Encapsulate(hostAddr)

			//记录地址以供接力
			ad := fullAddr.String()
			myaddr = ad
			log.Printf("我的fullAddr是 %s\n", fullAddr)

		}
	}
	return basicHost
}

//流处理
func handleStream(s net.Stream) {

	// 被别人连接时建立的流
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// 流s会持续打开直到一方关闭连接
}

func readData(rw *bufio.ReadWriter) {

	var str string
	for {
		str, _ = rw.ReadString(23)

		//读不到拉倒
		if str == "" {
			return
		}
		//读到了
		if str != string(23) {
			recv <- str
		}
	}
}

func writeData(rw *bufio.ReadWriter) {

	var sendchan = make(chan string, 100)
	SendList.PushBack(&sendchan)

	var str string
	for {
		str = <-sendchan
		p2pmutex.Lock()
		rw.WriteString(str)
		rw.Flush()
		p2pmutex.Unlock()

	}

}

func P2PSend(str string) {

	if strings.Contains(str, string(17)+string(18)) {
		return
	}
	str = strings.Replace(str, string(23), string(17)+string(18), -1)
	str = str + string(23)

	for p := SendList.Front(); p != nil; p = p.Next() {
		sc, ok := (p.Value).(*chan string)
		if ok {
			*sc <- str
		}
	}
	recv <- str
}

func P2PRecv() string {
	str := <-recv

	str = strings.TrimRight(str, string(23))
	str = strings.Replace(str, string(17)+string(18), string(23), -1)

	return str
}

func connectTarget(target string, ha host.Host) {
	// 从target参数解析ipfs地址
	ipfsaddr, _ := ma.NewMultiaddr(target)
	pid, _ := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	peerid, _ := peer.IDB58Decode(pid)

	// 解析PeerID和TargetID
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))

	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	// 加入节点列表
	ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

	//连接别人时建立的流
	s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
	if err != nil {
		fmt.Println("连接失败")
		fmt.Println(err)
		return
	}
	fmt.Println("已连接")

	// 建立流
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// 开始读写
	go writeData(rw)
	go readData(rw)

}

func getIDFromServer() int {
	conn, _ := csnet.Dial("tcp", *ServerAddrFlag+":"+P2PIDPort)
	defer conn.Close()
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, conn)
	st := strings.Replace(buf.String(), "\n", "", -1)
	id, _ := strconv.Atoi(st)
	return id
}

func updateNodeListFromServer() {
	for {
		CurrentNodelist, _ := json.Marshal(NodeList)
		CurrentNLstr := string(CurrentNodelist)
		time.Sleep(1e9)
		for {
			AddrBookConn, err := csnet.Dial("tcp", *ServerAddrFlag+":"+P2PAddrBookPort)
			if err != nil {
				log.Print("向Server请求更新NL时出错，Server可能已经下线")
				log.Print(err)
				return
			}
			Str, err := bufio.NewReader(AddrBookConn).ReadString('\n')
			Str = strings.Replace(Str, "\n", "", -1)
			if err != nil {
				log.Print("请求更新NL读取数据时出错")
				log.Print(err)
				return
			}
			AddrBookConn.Close()

			if Str != CurrentNLstr {
				log.Print("需要更新NodeList")
				_ = json.Unmarshal([]byte(Str), &NodeList)
				log.Print("更新NodeList成功")
				PrintNodeList()
				break
			}
			time.Sleep(1e9)
		}
	}
}

func PrintNodeList() {
	for _, x := range NodeList {
		fmt.Print("  ID")
		fmt.Print(x.ID)
		fmt.Print("  在线情况")
		fmt.Println(x.IsOnline)
	}
}

func GetAndSendAddr() {
	AddrBookConn, _ := csnet.Dial("tcp", *ServerAddrFlag+":"+P2PAddrBookPort)
	Str, _ := bufio.NewReader(AddrBookConn).ReadString('\n')
	fmt.Fprintf(AddrBookConn, strconv.Itoa(IDname)+"|")
	fmt.Fprintf(AddrBookConn, myaddr+"|")
	fmt.Fprint(AddrBookConn, string(PubKey))
	fmt.Fprint(AddrBookConn, string(24))
	//fmt.Fprint(AddrBookConn, "|")

	AddrBookConn.Close()
	_ = json.Unmarshal([]byte(Str), &NodeList)
	//spew.Dump(NodeList)
}

func P2PConnect() {

	ServerIP = flag.String("s", "", "连接IP")
	flag.Parse()

	SendList = list.New()

	var ha host.Host

	IDname = getIDFromServer()

	fmt.Print("我是第")
	fmt.Print(IDname)
	fmt.Println("号")

	RSAKeyGenInit()

	ha = makeBasicHost(15000 + IDname)
	ha.SetStreamHandler("/p2p/1.0.0", handleStream)

	GetAndSendAddr()

	if IDname != 1 { //如果要去连接其他机器

		for i := 1; i < len(NodeList); i++ {
			if NodeList[i].IsOnline == true {
				fmt.Printf("节点%d", NodeList[i].ID)
				connectTarget(NodeList[i].Addr, ha)

				PubKeyMap[strconv.Itoa(NodeList[i].ID)] = string(NodeList[i].PubKey)
				//spew.Dump(NodeList[i].PubKey)
			}
		}
		//写上自己的公钥
		PubKeyMap[strconv.Itoa(IDname)] = string(PubKey)
	}
}
