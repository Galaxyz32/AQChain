package main

import (
	"AQChain/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

var workspace = ""

func InsertMNDB(M models.MerkleNode, O orm.Ormer) { //插入MN一条记录
	_, err := O.Insert(&M)
	if err != nil {
		if DBLogPrint == true {
			beego.Info("插入MN失败插入MN失败插入MN失败插入MN失败插入MN失败", err)
		}

		return
	}
}
func InsertUserDB(U models.User, O orm.Ormer) { //插入user记录
	_, err := O.Insert(&U)
	if err != nil {
		if DBLogPrint == true {
			beego.Info("插入user失败插入user失败插入user失败插入user失败插入user失败", err)
		}

		return
	}
}
func InsertBlockDB(B models.Block, O orm.Ormer) { //插入user记录
	_, err := O.Insert(&B)
	if err != nil {
		if DBLogPrint == true {
			beego.Info("插入user失败插入user失败插入user失败插入user失败插入user失败", err)
		}

		return
	}
}

func InsertFile(F models.FileInBlockchain, O orm.Ormer) {
	_, err = O.Insert(&F)
	if err != nil {
		if DBLogPrint == true {
			beego.Info("插入文件失败插入文件失败插入文件失败插入文件失败", err)
		}

		return
	}
}

func ReadBlockFromDB(BBBB models.Block, OOOO orm.Ormer) {
	//time.Sleep(5e9)
	workspace = CaculatePath()
	var i int
	i = 1
	for {

		BBBB.Index = i

		err1 := OOOO.Read(&BBBB, "uid")
		if err1 != nil {
			fmt.Printf(" \x1b[32m本次共恢复%d个块 \x1b[0m", i-1)
			return
		} else {
			if DBLogPrint == true {
				fmt.Printf(" \x1b[32m找到序号为%d的块 \x1b[0m", BBBB.Index)
				fmt.Printf(" \x1b[32mContentHash是%s \x1b[0m", BBBB.Hash)

				Blockchain = append(Blockchain, BBBB)
			}
			i++
		}

	}

}

func DBSelect(MMMM models.MerkleNode, OOOO orm.Ormer) {
	workspace = CaculatePath()
	var i int
	i = 1
	for i < 10000 {
		MMMM.ContentHash, err = SHA256File(workspace + "/tmp/1.txt")
		//MMMM.ContentHash,err=SHA256File("./1.txt")
		//beego.Info("打印哈希值打印哈希值打印哈希值",MMMM.ContentHash)
		err1 := OOOO.Read(&MMMM, "ContentHash")
		if err1 != nil {
			if DBLogPrint == true {
				beego.Info("目标文件在Merkelnode库还未查询到，请稍等...", err1)
			}

			time.Sleep(10e9)
			i++
		} else {
			if MMMM.Blocktag == false {
				//beego.Info("", MMMM.ID)
				if DBLogPrint == true {
					fmt.Println("恭喜恭喜，目标文件在Merkelnode库查询成功！")
					fmt.Println("ID是：", MMMM.ID)
					fmt.Println("Blocktag是：", MMMM.Blocktag)
					fmt.Println("ContentHash是：", MMMM.ContentHash)
					fmt.Println("UserName是：", MMMM.UserName)
					fmt.Println("TimeStampOfContent是：", MMMM.TimeStampOfContent)
					fmt.Println("TypeOfThis是：", MMMM.TypeOfThis)
					fmt.Println("User2name是：", MMMM.User2name)
					fmt.Println("ValueOfMerkleNode是：", MMMM.ValueOfMerkleNode)
				}

				break
			}
		}
	}
}

func Fileselect(FFFF models.FileInBlockchain, OOOO orm.Ormer) {
	workspace = CaculatePath()
	var i int
	i = 1
	for i < 10000 {
		FFFF.FileContent, err = SHA256File(workspace + "2.txt")
		//MMMM.ContentHash,err=SHA256File("./1.txt")
		//beego.Info("打印哈希值打印哈希值打印哈希值",MMMM.ContentHash)
		err1 := OOOO.Read(&FFFF, "FileContent")
		if err1 != nil {
			if DBLogPrint == true {
				beego.Info("目标文件在文件库还未查询到，请稍等...", err1)
			}

			time.Sleep(10e9)
			i++
		} else {

			if DBLogPrint == true {
				fmt.Println("祝贺祝贺，目标文件在文件库查询成功！")
				fmt.Println("文件index是：", FFFF.Index)
				fmt.Println("文件creater是：", FFFF.Creater)
				fmt.Println("文件owner是：", FFFF.Owner)
				fmt.Println("文件filecontent是：", FFFF.FileContent)
				fmt.Println("文件value是：", FFFF.Value)
			}
			//beego.Info("", MMMM.ID)

			break
		}
	}
}
