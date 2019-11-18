package main

import (
	"fmt"
	"net/http"
)

func web(w http.ResponseWriter, r *http.Request) {

	//打印请求的方法
	//fmt.Println("method", r.Method)

	//否则走打印输出post接受的参数username和password
	status := r.PostFormValue("feedback")
	hash := r.PostFormValue("hash")
	OpenID := r.PostFormValue("OpenID")
	BlockIndex := r.PostFormValue("BlockIndex")
	MNIndex := r.PostFormValue("MNIndex")
	BlockCreater := r.PostFormValue("BlockCreater")
	BlockHash := r.PostFormValue("BlockHash")
	TimeStamp := r.PostFormValue("TimeStamp")

	//如果已接收到自动返回状态位200

	fmt.Println("文件" + hash)
	fmt.Println("OpenID" + OpenID)
	fmt.Println("状态" + status)
	fmt.Println("区块号" + BlockIndex)
	fmt.Println("块内交易号" + MNIndex)
	fmt.Println("区块创建者" + BlockCreater)
	fmt.Println("区块哈希" + BlockHash)
	fmt.Println("区块时间" + TimeStamp)

	fmt.Println()

}

func main() {
	http.HandleFunc("/web", web)
	http.ListenAndServe(":10200", nil)
	select {}
}
