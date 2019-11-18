package main

import (
	"AQChain/WJYBLOCK/impl"
	"AQChain/models"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var ( //转换器//薛
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { //允许跨域
			return true
		},
	}
	wsConn         *websocket.Conn
	err            error
	conn1          *impl.Connection
	M1             *models.MerkleNode
	jsonbyte1      []byte //薛JS中用
	WSmutex        = &sync.Mutex{}
	HtmlMNListener chan models.MerkleNode
)

//////////////////////////http//////////////////
func wshandler(w http.ResponseWriter, r *http.Request) { //hundlerFunc的回调函数
	//defer conn1.Close()

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		fmt.Println("case1")
		log.Print(err)
		return
	}
	if conn1, err = impl.InitConnection(wsConn); err != nil { //方法的调用为包名.方法名
		fmt.Println("case2")
		log.Print(err)
		return
	}

	var data []byte

	for {
		if data, err = conn1.ReadMessage(); err != nil {
			fmt.Println("ws读取时出错")
			//log.Print(err)
			return
		}
		if data != nil { //将前端传来的data放入message中，然后data清空

			switch string(data) {

			case "ZaiXian":
				WSmutex.Lock()
				log.Print("前端打开在线节点页面")
				OnlineNodesPage()
				WSmutex.Unlock()

			case "JiFen":
				WSmutex.Lock()
				log.Print("前端打开积分节点页面")

				WSmutex.Unlock()

			case "ShangChuan": //上传页面
				WSmutex.Lock()
				log.Print("前端打开上传页面")
				UploadPage()
				WSmutex.Unlock()

			case "ShangChuan2":
				WSmutex.Lock()
				fmt.Println("情况3")

				//data2, _ := conn1.ReadMessage()
				//for data2 != nil {
				//	fmt.Println(string(data2))
				//	data2, _ = conn1.ReadMessage()
				//	//if data2 == nil {
				//	//	break
				//	//}
				//	//HandlerMessageHTML(data2)
				//}
				WSmutex.Unlock()

			default:
				fmt.Println("前端返回类型不正确")
				fmt.Println(string(data))

			}

		}

	}
}

func BlockchainPage() {
	type BlockDetails struct{ index, creater, time, hash, prevHash string }
	var BlockList []BlockDetails
	for _, x := range Blockchain {
		bc := new(BlockDetails)
		bc.index = strconv.Itoa(x.Index)
		bc.creater = x.BlockCreater
		bc.time = x.TimeStamp
		bc.hash = x.Hash
		bc.prevHash = x.PrevBlockHash
		BlockList = append(BlockList, *bc)
	}
	spew.Dump(BlockList)
}

func MerkleNodePage() {
	//err := conn1.WriteMassage([]byte("Ready"))
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	//data, _ := conn1.ReadMessage()
	//MNNumber, _ := strconv.Atoi(string(data))
	MNNumber := 1
	if MNNumber >= len(Blockchain) {
		return
	}
	spew.Dump(Blockchain[MNNumber].BlockBody.BlockBody)
	//jsonbyte1,_:=json.Marshal(Blockchain[MNNumber].BlockBody.BlockBody)
}

//文件上传页后端处理逻辑
func UploadPage() {
	defer func() {
		fmt.Println("正在销毁通道")
		HtmlMNListener = nil
	}()

	//向前端发送数据
	go func() {
		fmt.Println("正在建立通道")
		HtmlMNListener = make(chan models.MerkleNode, 100)
		for {
			jsonbyte1, _ = json.Marshal(<-HtmlMNListener)
			conn1.WriteMassage(jsonbyte1)
		}
	}()

	//接收前端数据
	for {
		data2, _ := conn1.ReadMessage()
		if data2 != nil {
			log.Print("成功处理一条消息")
			HandlerMessageHTML(data2)
		} else {
			log.Print("退出上传页面消息处理")
			return
		}
	}
}
func HandlerMessageHTML(message []byte) {
	var input string
	//s := generateMerkleNode(string(message), 0, username, "", "")

	//鉴于网页上传暂时获取不到真实路径，暂改用随机生成mn方法
	path := generateFile(1)
	s := generateMerkleNode(path, 0, username, "WebUser", "")

	MNtoFileListenChan(s)
	input = string(SerializeMerkleNode(s))
	if len(input) > 0 {
		msg := input
		msg = "mn" + msg
		P2PSend(msg)
	}

}

func OnlineNodesPage() {
	type OnlineStatus struct{ Online, Offline int }
	resp := new(OnlineStatus)
	Online := 0
	Offline := -1

	for _, x := range NodeList {
		if x.IsOnline == true {
			Online++
		} else {
			Offline++
		}
	}

	resp.Offline = Offline
	resp.Online = Online

	spew.Dump(resp)

	jsonbyte1, err = json.Marshal(resp)
	if err != nil {
		log.Print(err)
	}
	conn1.WriteMassage(jsonbyte1)

}

func ContributionPage() {
	id, _ := strconv.Atoi(username)
	MyContribution := UserTable[id-1].Contribution
	fmt.Println(MyContribution)
	//jsonbyte1, _ = json.Marshal(resp)
	//conn1.WriteMassage(jsonbyte1)

}

var MyFileList []ProtectedFile

type ProtectedFile struct {
	name, time, hash string
	BlockGenerated   bool
}

func MNtoFileListenChan(node models.MerkleNode) {
	MyFile := new(ProtectedFile)
	MyFile.hash = node.ContentHash
	MyFile.name = node.User2name
	MyFile.time = node.TimeStampOfContent
	MyFile.BlockGenerated = false
	MyFileList = append(MyFileList, *MyFile)
}

func MNVerifiedbyBlock(block models.Block) {
	for _, x := range block.BlockBody.BlockBody {
		for i, y := range MyFileList {
			if y.hash == x.ContentHash {
				MyFileList[i].BlockGenerated = true
			}
		}
	}
}

func ProtectedFilePage() {
	jsonbyte1, _ = json.Marshal(MyFileList)
	conn1.WriteMassage(jsonbyte1)
}
