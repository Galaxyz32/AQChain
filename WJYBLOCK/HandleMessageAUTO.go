package main

import (
	"AQChain/models"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

func HandlerMessageAUTO(FFile models.FileInBlockchain) {

	var input string
	var path string
	var Filename int64
	Filename = 1
	var o orm.Ormer

	//HMST := getIntFromEnv("HandleMessageSleepTimeSec")
	time.Sleep(time.Second * time.Duration(getIntFromEnv("InitWaitTimeSec")))
	for username != "" {
		mnuser1 := RandInt(1, NODESNUM+1) //表示这一次该节点需不需要操作
		mnuser2 := fmt.Sprint(mnuser1)
		//action1 := getIntFromEnv("Action1Num")
		//action2 := getIntFromEnv("Action2Num")

		if username == mnuser2 {
			//if mnuser1 >= 1 && mnuser1 <= action1 { //表明是知识产权上传者
			path = generateFile(Filename)

			s := generateMerkleNode(path, 0, username, "", "")
			input = string(SerializeMerkleNode(s))
			for flag1 == false {
				if SilentMode == false {
					fmt.Println("现在正在生成block，需要等待")
				}
				time.Sleep(1e9)
				//time.Sleep(time.Second * time.Duration(HMST)) //2s表示现在block生成大约需要多久完成
			}

			if DebugNoDB == false {
				////////////////////////////////插入文件表相关////////////////////////////////////

				FFile.Index = Filename
				FFile.Creater = username
				FFile.Owner = username //不交易，与上面相同
				FFile.FileContent, err = SHA256File(path)
				FFile.Value = 0 //不交易，价值为零
				////////////////////////////////插入文件表相关////////////////////////////////////
				o = OO
				InsertFile(FFile, o)
			}

			if len(input) > 0 {
				msg := input
				msg = "mn" + msg
				P2PSend(msg)
			}

			time.Sleep(1e9)
			Filename++
			//} else if mnuser1 > action1 && mnuser1 <= action2 { //表示交易者
			//	//随机从file表中找到一个文件，提取创作者、文件名和价值
			//	if len(FileTable) != 0 {
			//		randindex := RandInt(1, len(FileTable)-1) //需要减1，才能不超边界
			//		flagfilemn = true                         //在本周期中，为了防止重复交易，设置flagfilemn
			//		file2 := FileTable[randindex]
			//		for i := 0; i < len(filemn); i++ {
			//			if file2.FileContent == filemn[i] {
			//				if SilentMode == false {
			//					fmt.Println("该文件已经被交易")
			//				}
			//				flagfilemn = false
			//			}
			//		}
			//		if flagfilemn {
			//			file2owner := file2.Creater
			//			file2value := file2.Value
			//			file2name := file2.FileContent //filecontent是文件的哈希值
			//			var randvalue float64
			//			if file2value == 0 {
			//				randvalue = float64(RandInt(1, 100))
			//			} else {
			//				randvalue = file2value + file2value*float64(RandInt(1, 8))
			//			}
			//			s := generateMerkleNode("", randvalue, file2owner, username, file2name)
			//			input = string(SerializeMerkleNode(s))
			//			for flag1 == false {
			//				if SilentMode == false {
			//					fmt.Println("现在正在生成block，需要等待")
			//				}
			//				time.Sleep(time.Second * time.Duration(HMST)) //2s表示现在block生成大约需要多久完成
			//			}
			//			time.Sleep(0.5e8)
			//			//	time.Sleep(time.Second * 5)
			//			////////////////////////////////更新文件表相关（因为交易）////////////////////////////////////
			//
			//			////////////////////////////////更新文件表相关（因为交易）////////////////////////////////////
			//			if len(input) > 0 {
			//				msg := input
			//				msg = "mn" + msg
			//				P2PSend(msg)
			//			}
			//		}
			//	}
			//} else if mnuser1 > action2 && mnuser1 <= NODESNUM {
			//	randcaozuo := RandInt(1, 2) //1表示上传知识产权，2表示交易
			//	if randcaozuo == 1 {
			//		//同1-10的上传
			//		path = generateFile(Filename)
			//
			//		s := generateMerkleNode(path, 0, username, "", "")
			//		input = string(SerializeMerkleNode(s))
			//		for flag1 == false {
			//			if SilentMode == false {
			//				fmt.Println("现在正在生成block，需要等待")
			//			}
			//			//			time.Sleep(time.Second * time.Duration(HMST)) //2s表示现在block生成大约需要多久完成
			//			time.Sleep(0.5e8) //2s表示现在block生成大约需要多久完成
			//
			//		}

			//if DebugNoDB == false {
			//	////////////////////////////////插入文件表相关////////////////////////////////////
			//	FFile.Index = Filename
			//	FFile.Creater = username
			//	FFile.Owner = username //不交易，与上面相同
			//	FFile.FileContent = hash(path)
			//	FFile.Value = 0 //不交易，价值为零
			//	////////////////////////////////插入文件表相关////////////////////////////////////
			//	o = OO
			//	InsertFile(FFile, o)
			//}

			//time.Sleep(1e9)
			////	time.Sleep(time.Second * 5)
			//if len(input) > 0 {
			//	msg := input
			//	msg = "mn" + msg
			//	P2PSend(msg)
			//}
			//} else if randcaozuo == 2 {
			//	//同11-20的交易
			//	//随机从file表中找到一个文件，提取创作者、文件名和价值
			//	//选择交易对象的时候需要注意是否在这一批中已经被别的节点选中，如果filetable中没有文件，则无法交易
			//	if len(FileTable) != 0 {
			//		randindex := RandInt(1, len(FileTable)-1)
			//		flagfilemn := true
			//		file2 := FileTable[randindex]
			//		for i := 0; i < len(filemn); i++ {
			//			if file2.FileContent == filemn[i] {
			//				if SilentMode == false {
			//					fmt.Println("该文件已经被交易")
			//				}
			//				flagfilemn = false
			//			}
			//		}
			//		if flagfilemn {
			//			file2owner := file2.Creater
			//			file2value := file2.Value
			//			file2name := file2.FileContent //filecontent是文件的哈希值
			//			randvalue := file2value + file2value*float64(RandInt(1, 8))
			//			s := generateMerkleNode("", randvalue, file2owner, username, file2name)
			//			input = string(SerializeMerkleNode(s))
			//			for flag1 == false {
			//				if SilentMode == false {
			//					fmt.Println("现在正在生成block，需要等待")
			//				}
			//				time.Sleep(time.Second * time.Duration(HMST)) //2s表示现在block生成大约需要多久完成
			//			}
			//			////////////////////////////////更新文件表相关（因为交易）////////////////////////////////////
			//
			//			////////////////////////////////更新文件表相关（因为交易）////////////////////////////////////
			//			time.Sleep(0.5e8)
			//			//		time.Sleep(time.Second * 5)
			//			if len(input) > 0 {
			//				msg := input
			//				msg = "mn" + msg
			//				P2PSend(msg)
			//			}
			//		}
			//	}
			//
			//}

		}

		//time.Sleep(time.Second * time.Duration(HMST) * 10)
		time.Sleep(1e9)

	}
}
