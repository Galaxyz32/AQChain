package main

import (
	"AQChain/models"
	"container/list"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WXConfirmMsg struct {
	BlockIndex   int
	MNIndex      int
	BlockCreater string
	BlockHash    string
	FileHash     string
	TimeStamp    string
}

var WXIPT = make(chan *WXInput, 100)
var WXOPT = make(chan *WXOutput, 100)
var WXFileConfirmList *list.List
var WXFileTimeOutList *list.List

type WXOutput struct {
	Hash   string
	Status int
}
type WXInput struct {
	OpenID string
	Hash   string
	Name   string
}

func WXFileMonitor(hash string, Type string, OpenID string) {

	ExistFlag := false
	for i := 0; i < len(Blockchain); i++ {
		for j := 0; j < len(Blockchain[i].BlockBody.BlockBody); j++ {
			if hash == Blockchain[i].BlockBody.BlockBody[j].ContentHash {
				ExistFlag = true
			}
		}
	}
	for i := 0; i < len(blockbody.BlockBody); i++ {
		if hash == blockbody.BlockBody[i].ContentHash {
			ExistFlag = true
		}
	}
	for t := 0; t < 30; t++ {
		for p := WXFileConfirmList.Front(); p != nil; p = p.Next() {
			if p.Value.(WXConfirmMsg).FileHash == hash {
				if ExistFlag == false {

					respon, err := http.PostForm("https://www.bjutxzpj.cn/protect/apiwx/Wxreceive/blockfeedback"+Type,
						//_, err := http.PostForm("http://127.0.0.1:10200/web",
						url.Values{"hash": {hash}, "feedback": {"2"},
							"OpenID":       {OpenID},
							"BlockIndex":   {strconv.Itoa(p.Value.(WXConfirmMsg).BlockIndex)},
							"MNIndex":      {strconv.Itoa(p.Value.(WXConfirmMsg).MNIndex + 1)},
							"BlockCreater": {p.Value.(WXConfirmMsg).BlockCreater},
							"BlockHash":    {p.Value.(WXConfirmMsg).BlockHash},
							"TimeStamp":    {p.Value.(WXConfirmMsg).TimeStamp},
						})
					fmt.Println("小程序返回状态码" + respon.Status)
					body, _ := ioutil.ReadAll(respon.Body)
					fmt.Print("小程序接收服务器返回值")
					fmt.Println(string(body))

					if err != nil {
						fmt.Println("向小程序服务器反馈上链成功失败")
					}

				} else {
					_, _ = http.PostForm("https://www.bjutxzpj.cn/protect/apiwx/Wxreceive/blockfeedback"+Type,
						//_, err := http.PostForm("http://127.0.0.1:10200/web",
						url.Values{"hash": {hash}, "feedback": {"3"},
							"OpenID":       {OpenID},
							"BlockIndex":   {strconv.Itoa(p.Value.(WXConfirmMsg).BlockIndex)},
							"MNIndex":      {strconv.Itoa(p.Value.(WXConfirmMsg).MNIndex + 1)},
							"BlockCreater": {p.Value.(WXConfirmMsg).BlockCreater},
							"BlockHash":    {p.Value.(WXConfirmMsg).BlockHash},
							"TimeStamp":    {p.Value.(WXConfirmMsg).TimeStamp},
						})
					if err != nil {
						fmt.Println("向小程序服务器反馈重复文件上链成功失败")
					}
				}
				WXFileConfirmList.Remove(p)
				return
			}
		}
		time.Sleep(20e9)
	}
	_, err := http.PostForm("https://www.bjutxzpj.cn/protect/apiwx/Wxreceive/blockfeedback"+Type,
		//_, err := http.PostForm("http://127.0.0.1:10200/web",
		url.Values{"hash": {hash}, "feedback": {"4"}, "OpenID": {OpenID}})
	if err != nil {
		fmt.Println("向小程序服务器反馈上链失败失败")
	}
	WXFileTimeOutList.PushBack(hash)

}

func WXTimeOutFileViewer() {
	for {
		spew.Dump(WXFileTimeOutList)
		fmt.Println()
		time.Sleep(10e9)
	}

}

func ListenWXServer() {
	http.HandleFunc("/wx", WXIPTHandler)
	http.HandleFunc("/wx/sh", WXIPTHandler)
	http.HandleFunc("/wx/bj", WXIPTHandler)
	http.HandleFunc("/wx/gz", WXIPTHandler)
	http.ListenAndServe(":18888", nil)
}

func WXIPTHandler(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	fmt.Println("接收微信消息来自" + r.Method)
	if r.PostFormValue("OpenID") == "" || r.PostFormValue("Name") == "" || r.PostFormValue("Hash") == "" {
		w.WriteHeader(202)
	}
	msg := new(WXInput)
	msg.OpenID = r.PostFormValue("OpenID")
	msg.Name = r.PostFormValue("Name")
	msg.Hash = r.PostFormValue("Hash")
	TypeOfContent := r.PostFormValue("Type")
	WXIPT <- msg
	go WXFileMonitor(msg.Hash, TypeOfContent, msg.OpenID)
	//如果已接收到自动返回状态位200
}

func HandlerMessageWX() {
	WXFileConfirmList = list.New()
	WXFileTimeOutList = list.New()

	var input string
	HMST := getIntFromEnv("HandleMessageSleepTimeSec")

	for username != "" {
		recvmsg := <-WXIPT

		newMN := models.MerkleNode{}
		newMN.UserName = username
		newMN.User2name = "wx// " + (*recvmsg).OpenID + " // " + (*recvmsg).Name
		newMN.ContentHash = (*recvmsg).Hash
		newMN.TimeStampOfContent = time.Now().Format(TimeType)
		newMN.TypeOfThis = 0
		newMN.ValueOfMerkleNode = 0

		input = string(SerializeMerkleNode(newMN))

		for flag1 == false {
			if SilentMode == false {

				fmt.Println("现在正在生成block，需要等待")
			}
			time.Sleep(time.Second * time.Duration(HMST)) //2s表示现在block生成大约需要多久完成
		}

		if len(input) > 0 {
			msg := input
			msg = "mn" + msg
			P2PSend(msg)
		}

	}

}
