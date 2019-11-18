package main

import (
	"fmt"
)

func CalculateContribution() {
	maxCIndex = 0

	lastBlock := Blockchain[len(Blockchain)-1]
	for i := 0; i < len(lastBlock.BlockBody.BlockBody); i++ {

		for j := 0; j < len(UserTable); j++ {
			if UserTable[j].UserName == lastBlock.BlockBody.BlockBody[i].UserName {
				UserTable[j].Contribution++
			}
		}
	}

	for j := 0; j < len(UserTable); j++ {

		if UserTable[j].Contribution > UserTable[maxCIndex].Contribution {
			maxCIndex = j
		}

	}

	//for i := 0; i < len(Blockchain[len(Blockchain)-1].BlockBody.BlockBody); i++ {
	//	if len(FileTable) == 0 { //当filetable中没有数据的时候，进行该操作
	//		if SilentMode == false {
	//
	//			fmt.Println("file1111111111111111")
	//		}
	//		filet := models.FileInBlockchain{}
	//		filet.Owner = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].User2name
	//		filet.Value = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ValueOfMerkleNode
	//		filet.FileContent = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ContentHash
	//		filet.Creater = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].UserName
	//		FileTable = append(FileTable, filet)
	//	} else {
	//		for j := 0; j < len(FileTable); j++ {
	//			if Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ContentHash == FileTable[j].FileContent && Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ValueOfMerkleNode != 0 {
	//				//表示最后一个区块中mn的文件哈西与文件表中某一文件哈希一致，并且最后一个区块的mnvalue不为0，代表这是一个针对该文件的交易。所以更新file表中相应的数据
	//				FileTable[j].Value = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ValueOfMerkleNode
	//				FileTable[j].Owner = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].User2name
	//				if SilentMode == false {
	//					fmt.Println("filefilefilefilefile1")
	//				}
	//
	//				break
	//			} else if Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ContentHash == FileTable[j].FileContent && Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ValueOfMerkleNode == 0 {
	//				//表示最后一个区块中mn的文件哈西与文件表中某一文件哈希一致，但mn中的文件价值为0，说明这是一笔错误数据。输出错误信息
	//				if SilentMode == false {
	//					fmt.Println("该文件已经存在于区块链中")
	//				}
	//
	//				break
	//			} else {
	//				if SilentMode == false {
	//
	//					fmt.Println("filefilefilefilefile2")
	//				}
	//				//其他的情况就是最后一个区块中的mn没有对应的文件表中的文件，属于新的，所以加入到文件表
	//				filet := models.FileInBlockchain{}
	//				filet.Owner = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].User2name
	//				filet.Value = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ValueOfMerkleNode
	//				filet.FileContent = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].ContentHash
	//				filet.Creater = Blockchain[len(Blockchain)-1].BlockBody.BlockBody[i].UserName
	//				FileTable = append(FileTable, filet)
	//				break
	//			}
	//		}
	//	}
	//
	//}

	//将节点的贡献值写入到excel
	//if *WebFlag == false && *WeiXinFlag == false {
	//	writeEXCEL()
	//}

	FileMutex.Lock()
	fmt.Fprintf(CtbWriter, "第%d个块: \n", CtbIndex)
	for i := 0; i < len(UserTable); i++ {
		fmt.Fprintf(CtbWriter, "%s : %.2f\n", UserTable[i].UserName, UserTable[i].Contribution)
	}
	fmt.Fprintln(CtbWriter)
	CtbWriter.Flush()
	FileMutex.Unlock()

	CtbIndex++

	UserTable[maxCIndex].Contribution = 0
	UserTable[maxCIndex].Contribution1 = 0
	UserTable[maxCIndex].Contribution2 = 0
	UserTable[maxCIndex].UserFileNume = 0

	////计算所占空间大小
	//if SilentMode == false {
	//	fmt.Println("FileTable size :", unsafe.Sizeof(FileTable))
	//	fmt.Println("Blockchain size :", unsafe.Sizeof(Blockchain))
	//	fmt.Println("UserTable size :", unsafe.Sizeof(UserTable))
	//}

}
