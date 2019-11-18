package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	"io"
	"io/ioutil"
	"log"
	csnet "net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var myaddr string
var IDname int
var filepath string

// 地址本
type NodeStruct struct {
	ID       int
	Addr     string
	IsOnline bool
}

var NodeList []NodeStruct

//读写锁
var mutex = &sync.Mutex{}

//指定监听
func makeBasicHost(listenPort int) host.Host {

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
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

	log.Println("收到传入连接!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// 流s会持续打开直到一方关闭连接
}

func readData(rw *bufio.ReadWriter) {
	if IDname == 2 {
		str, _ := rw.ReadString(23)
		str = strings.TrimRight(str, string(23))
		str = strings.Replace(str, string(17)+string(18)+string(19)+string(20), string(23), -1)
		spew.Dump([]byte(str))

		outputFile, outputError := os.OpenFile(filepath+"2.mp3", os.O_WRONLY|os.O_CREATE, 0666)
		if outputError != nil {
			fmt.Printf("An error occurred with file opening or creation\n")
			return
		}
		outputFile.Write([]byte(str))
		outputFile.Close()

		fmt.Println("接收完成")
	}

}

func writeData(rw *bufio.ReadWriter) {
	if IDname == 1 {
		file, err := os.Open(filepath + "1.mp3")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		contents, _ := ioutil.ReadAll(file)
		spew.Dump(contents)

		str := string(contents)
		str = strings.Replace(str, string(23), string(17)+string(18)+string(19)+string(20), -1)

		rw.WriteString(str)
		rw.WriteByte(23)
		rw.Flush()

		fmt.Println("发送完成")
	}
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

	s, _ := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")

	// 建立流
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// 开始读写
	go writeData(rw)
	go readData(rw)

}

func getID() int {
	conn, _ := csnet.Dial("tcp", ":19099")
	defer conn.Close()
	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, conn)
	st := strings.Replace(buf.String(), "\n", "", -1)
	id, _ := strconv.Atoi(st)
	return id
}

func GetAndSendAddr() {
	AddrBookConn, _ := csnet.Dial("tcp", ":19093")
	Str, _ := bufio.NewReader(AddrBookConn).ReadString('\n')
	fmt.Fprintf(AddrBookConn, myaddr+"\n")
	AddrBookConn.Close()
	_ = json.Unmarshal([]byte(Str), &NodeList)
}

func main() {
	var ha host.Host

	filepath = CaculatePath()

	IDname = getID()
	fmt.Println(IDname)

	ha = makeBasicHost(15000 + IDname)
	ha.SetStreamHandler("/p2p/1.0.0", handleStream)

	GetAndSendAddr()

	if IDname != 1 {

		for i := 1; i < len(NodeList); i++ {
			if NodeList[i].IsOnline == true {
				connectTarget(NodeList[i].Addr, ha)
				fmt.Printf("节点%d已连接\n", NodeList[i].ID)
			}
		}

	}

	fmt.Println("_________等待传入___________")

	select {} // 阻塞

}
