package main

import (
	"AQChain/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var blockbody models.BlockBody

func OpenTheWeb(Mmm models.MerkleNode, Ooo orm.Ormer, Uss models.User, Bll models.Block, Fll models.FileInBlockchain) {

	var BlockValidTemp []models.MerkleNode
	BlockValidList := make(chan string, 100)

	flag1 = true //true表示可以上传mn，当前没有在生成block。因为生成block后会清空blockbody，
	// 所有为了防止中间有进来的mn被清空加入此操作

	if *WebFlag == true {
		go func() {
			http.HandleFunc("/ws/set", wshandler) //配置路由,ws代表websocket
			http.ListenAndServe("127.0.0.1:7777", nil)
		}()
	}

	username = strconv.Itoa(IDname)
	fmt.Println("我是第" + username + "号")
	//if len(Blockchain)==0 {
	//	b := gerateGensisBlocks()
	//	if *VerboseFlag == false {
	//		fmt.Println("创世区块")
	//		spew.Dump(b)
	//	}
	//}

	//
	//fmt.Printf("\x1b[34m")
	//spew.Dump(Blockchain)
	//fmt.Printf("\x1b[0m")

	//发送自己的ip地址到其他客户端
	//user := models.User{}
	Uss.UserName = username
	Uss.UserAddTime = time.Now().Format(TimeType)
	id := "id" + string(SerializeUser(Uss))
	//conn.Write([]byte(id))
	P2PSend(id)

	if *WebFlag == false && *WeiXinFlag == false {
		go HandlerMessageAUTO(Fll)
		////////////////////////////////插入文件表////////////////////////////////////
		//在HandlerMessageAUTO中进行文件表内容的插入
		////////////////////////////////插入文件表////////////////////////////////////
	}

	if DebugNoDB == false {
		//////////////////////////User插入数据库///////////////////////////////////////////
		InsertUserDB(Uss, Ooo)
		//////////////////////////User插入数据库///////////////////////////////////////////
	}

	if *WeiXinFlag == true {
		go HandlerMessageWX()
	}

	time.Sleep(1e9)
	//接受来自服务器的消息
	blockbody := models.BlockBody{}
	maxCIndex = 0

	var str string
	for {
		str = P2PRecv()

		var Type = ""
		MsgHash := hash(str)
		if strings.Contains(str, "id") {
			Type += "I "
			fmt.Fprintln(RecvWriter, Type+MsgHash)
		}
		if strings.Contains(str, "mn") {
			Type += "M "
			fmt.Fprintln(RecvWriter, Type)
		}
		if strings.Contains(str, "block11") {
			Type += "B "
			fmt.Fprintln(RecvWriter, Type+MsgHash)
		}
		if strings.Contains(str, "DBcontent") {
			Type += "D "
			fmt.Fprintln(RecvWriter, Type+MsgHash)
		}
		RecvWriter.Flush()

		if strings.Contains(str, "id") {
			s := str
			s = strings.Trim(s, "id")
			u := DeserializeUser([]byte(s))

			DebugInfo := "新用户 ID " + u.UserName + " 时间 " + u.UserAddTime
			fmt.Fprintln(MsgWriter, DebugInfo)
			//fmt.Fprintln(MsgWriter)
			MsgWriter.Flush()

			if *VerboseFlag == false {
				spew.Dump(u)
			}
			if username == "1" && len(UserTable) == 0 { //第一个节点，用户表为空
				UserTable = append(UserTable, Uss) //将自己的信息加入用户表
				if SilentMode == false {

					fmt.Println("我是第一个用户，我正在写入")
					//writeNewUserToDB(0)

				}

				usertablebytes, _ := json.Marshal(UserTable)
				P2PSend("DBcontent" + string(usertablebytes))
			} else if username == "1" && len(UserTable) != 0 { //目前的client是1，新加入的id是2，
				UserTable = append(UserTable, *u)
				usertablebytes, _ := json.Marshal(UserTable)
				P2PSend("DBcontent" + string(usertablebytes))
			} else if len(UserTable) != 0 { //网络中已有节点除第一个节点外将新节点加入本地数据库
				UserTable = append(UserTable, *u)
				if *VerboseFlag == false {
					//spew.Dump(UserTable)
				}

			}

		} else if strings.Contains(str, "mn") {
			s := str
			s = strings.Trim(s, "mn")

			Mmm = *DeserializeMerkleNode([]byte(s))

			DebugInfo := "新MN ID " + Mmm.UserName + " 类型 " + strconv.Itoa(Mmm.TypeOfThis) + " 时间 " + Mmm.TimeStampOfContent
			fmt.Fprintln(MsgWriter, DebugInfo)
			//fmt.Fprintln(MsgWriter)
			MsgWriter.Flush()

			if *VerboseFlag == false {
				spew.Dump(Mmm)
			}

			//fmt.Printf("ID: %d, ERR: %v\n", id, err)
			blockbody.BlockBody = append(blockbody.BlockBody, Mmm)

			if DebugNoDB == false {
				///////////////////////////////////////MerkleNode插入数据库///////////////////////////////////
				InsertMNDB(Mmm, Ooo)
				///////////////////////////////////////MerkleNode插入数据库///////////////////////////////////

			}

			if *WebFlag == true && HtmlMNListener == nil {
				log.Print("此时还没打开网页")
			}

			if *WebFlag == true && HtmlMNListener != nil {
				log.Print("已经打开网页")
				HtmlMNListener <- Mmm
			}

			filemn = append(filemn, Mmm.ContentHash)
			if SilentMode == false {
				fmt.Println("aaaaaaaaaaaaaaaaa")
				fmt.Printf(" \nMerkerNode	%d\n ", len(blockbody.BlockBody))
			}

			if len(blockbody.BlockBody)%MNPerBlock == 0 {
				if SilentMode == false {
					fmt.Println("是时候生成真正的区块了！")
				}

				sort.Sort(models.SortByHash{blockbody})
				//fmt.Fprintln(MsgWriter,"缓存区mn情况")
				//for i:=0;i<len(blockbody.BlockBody) ;i++  {
				//	fmt.Fprintln(MsgWriter,blockbody.BlockBody[i].UserName+" "+ blockbody.BlockBody[i].TimeStampOfContent)
				//}
				////fmt.Fprintln(MsgWriter)
				//MsgWriter.Flush()
				BlockValidTemp = blockbody.BlockBody
				BlockValidList <- SHA256Byte(SerializeBlockBody(blockbody.BlockBody))

				flag1 = false //false表示正在生成新区块，所有节点上传mn需要等到mn为true时再进行，但由于这样需要时刻计算flag1的值，所以复杂一些，在测试中使用往后推迟2s的方式来处理。

				if username == UserTable[maxCIndex].UserName { //如果应该生成新区块的是我
					newblock := generateNewBlock(blockbody, UserTable[maxCIndex].UserName)
					Bll = newblock
					if SilentMode == false {
						fmt.Println("正在生成新区块 ")
					}

					if DebugNoDB == false {
						////////////////////////////////////Block插入数据库//////////////////////////////
						//InsertBlockDB(Bll, Ooo)
						////////////////////////////////////Block插入数据库//////////////////////////////
					}

					a := SerializeBlock(newblock)
					b := "block11" + string(a)

					//time.Sleep(0.5e9)

					P2PSend(b)

				}
			}
		} else if strings.Contains(str, "block11") {
			s := str
			s = strings.Trim(s, "block11")
			b := DeserializeBlock([]byte(s))

			fmt.Fprintln(MsgWriter)
			DebugInfo := "新区块 ID " + b.BlockCreater + " 区块序号 " + strconv.Itoa(b.Index) + " 时间 " + b.TimeStamp
			fmt.Fprintln(MsgWriter)

			fmt.Fprintln(MsgWriter, DebugInfo)
			//fmt.Fprintln(MsgWriter)

			sort.Sort(models.SortByHash{b.BlockBody})
			//fmt.Fprintln(MsgWriter,"接受块mn情况")
			//for i:=0;i<len(b.BlockBody.BlockBody) ;i++  {
			//	fmt.Fprintln(MsgWriter, b.BlockBody.BlockBody[i].UserName+" "+b.BlockBody.BlockBody[i].TimeStampOfContent)
			//}
			//MsgWriter.Flush()

			//fmt.Printf("\x1b[32m")
			if SHA256Byte(SerializeBlockBody(b.BlockBody.BlockBody)) == <-BlockValidList {
				fmt.Printf("\x1b[32m  块校验通过  \x1b[0m \n")
				//LogMsg:="OK  由 "+b.BlockCreater+" 发送的块 "+strconv.Itoa(b.Index)+" 校验通过"
				//fmt.Fprintln(LogWriter, LogMsg)
				//LogWriter.Flush()
			} else {
				fmt.Printf("\x1b[31m  块校验失败  \x1b[0m \n")
				LogMsg := "块  " + b.BlockCreater + " 发送的块 " + strconv.Itoa(b.Index) + " 校验出错"
				fmt.Fprintln(LogWriter, LogMsg)
				fmt.Fprintln(LogWriter, "缓存区数据")
				fmt.Fprintln(LogWriter, spew.Sdump(BlockValidTemp))

				fmt.Fprintln(LogWriter, "区块数据")
				fmt.Fprintln(LogWriter, spew.Sdump(b.BlockBody))

				LogWriter.Flush()

			}

			//fmt.Printf("\x1b[32m")
			if b.BlockCreater == UserTable[maxCIndex].UserName {
				fmt.Printf("\x1b[32m  发送者校验通过  \x1b[0m \n")
				//LogMsg:="OK  由 "+b.BlockCreater+" 发送的块 "+strconv.Itoa(b.Index)+" 校验通过"
				//fmt.Fprintln(LogWriter, LogMsg)
				//LogWriter.Flush()
			} else {
				fmt.Printf("\x1b[31m  发送者校验失败  \x1b[0m \n")
				LogMsg := "发  " + b.BlockCreater + " 发送的块 " + strconv.Itoa(b.Index) + " 应由" + UserTable[maxCIndex].UserName + "发送"
				fmt.Fprintln(LogWriter, LogMsg)
				LogWriter.Flush()
			}

			if b.PrevBlockHash == Blockchain[len(Blockchain)-1].Hash {
				fmt.Printf("\x1b[32m  哈希值校验通过  \x1b[0m \n")
				//LogMsg:="序  "+b.BlockCreater+" 发送的块 "+strconv.Itoa(b.Index)+"Hash校验通过"
				//fmt.Fprintln(LogWriter, LogMsg)
				//LogWriter.Flush()
			} else {
				fmt.Printf("\x1b[31m  哈希值校验失败  \x1b[0m \n")
				LogMsg := "序  " + b.BlockCreater + " 发送的块 " + strconv.Itoa(b.Index) + "前序哈希不对"
				fmt.Fprintln(LogWriter, LogMsg)
				LogWriter.Flush()
			}

			Blockchain = append(Blockchain, *b)
			MNVerifiedbyBlock(*b)
			if DebugNoDB == false {
				////////////////////////////////////Block插入数据库//////////////////////////////
				InsertBlockDB(*b, Ooo)
				////////////////////////////////////Block插入数据库//////////////////////////////
			}

			if *WeiXinFlag == true {
				for i := 0; i < len(b.BlockBody.BlockBody); i++ {
					if strings.HasPrefix(b.BlockBody.BlockBody[i].User2name, "wx") {

						WXFeedBack := WXConfirmMsg{b.Index, i, b.BlockCreater,
							b.Hash, b.BlockBody.BlockBody[i].ContentHash, b.TimeStamp}
						//WXFileConfirmList.PushBack(b.BlockBody.BlockBody[i].ContentHash)
						WXFileConfirmList.PushBack(WXFeedBack)

					}
				}
			}

			//spew.Dump(b)

			//writeNewBlockToDB(len(Blockchain)-1)

			SliceClear2(&blockbody.BlockBody)
			SliceClear3(&filemn)
			flag1 = true //true表示可以继续上传mn
			go CalculateContribution()
		} else if strings.Contains(str, "DBcontent") { //利用DBcontent对现有的用户表进行覆盖
			s := str
			s = strings.Trim(s, "DBcontent")
			//fmt.Println("收到了一个DBcontent")
			DBcontentbuf := []byte(s)
			//fmt.Println("111111111111111111111111111111")
			if PathExists(DBPath) == false {
				json.Unmarshal(DBcontentbuf, &UserTable)

				DebugInfo := "新用户表内含有"
				fmt.Fprintln(MsgWriter, DebugInfo)

				for i := 0; i < len(UserTable); i++ {
					fmt.Fprintln(MsgWriter, UserTable[i].UserName+" "+UserTable[i].UserAddTime)
				}
				//fmt.Fprintln(MsgWriter)

				MsgWriter.Flush()

				//if *VerboseFlag == false {
				//	fmt.Println("现在的Usertable是")
				//	spew.Dump(UserTable)
				//}

				//fmt.Println("写入数据库")
				if *VerboseFlag == false {
					for i := 0; i < len(UserTable); i++ {
						spew.Dump(UserTable[i])
						//writeNewUserToDB(i)
					}
				}
			}
		}
	}
}
