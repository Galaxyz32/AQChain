package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/tealeg/xlsx/xlsx"
	"os"
	"strconv"
	"time"
)

func hash(str string) (s string) {

	//使用sha256哈希函数
	h := sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	//由于是十六进制表示，因此需要转换
	return hex.EncodeToString(sum)

}

func generateFile(i int64) string {
	filename := strconv.FormatInt(i, 10)
	file6, _ := os.Create(WorkSpace + "/tmp/" + filename + ".txt")
	data := strconv.FormatInt(time.Now().UnixNano(), 10)
	file6.Write([]byte(data))
	file6.Close()
	return WorkSpace + "/tmp/" + filename + ".txt"
}

func getIntFromEnv(key string) int {
	text := os.Getenv(key)
	num, _ := strconv.Atoi(text)
	return num
}

func writeEXCEL() {
	if username == "1" { //只需要一个节点来记录贡献值的变化，让第一个节点记录
		for i := 0; i < len(UserTable); i++ {
			//fmt.Println(1111111111111 + len(UserTable))
			//spew.Dump(UserTable)
			indexExcel, _ := strconv.Atoi(UserTable[i].UserName)
			valueExcel := strconv.FormatFloat(Decimal(UserTable[i].Contribution), 'f', -1, 64)
			//fmt.Println(indexExcel)
			//fmt.Println(valueExcel)
			AddToExcel(WorkSpace+"test1c.xlsx", indexExcel, valueExcel)
		}
		for i := 0; i < len(UserTable); i++ {
			indexExcel, _ := strconv.Atoi(UserTable[i].UserName)
			valueExcel := strconv.FormatFloat(UserTable[i].Contribution1, 'f', -1, 64)
			AddToExcel(WorkSpace+"test1c1.xlsx", indexExcel, valueExcel)
		}
		for i := 0; i < len(UserTable); i++ {
			indexExcel, _ := strconv.Atoi(UserTable[i].UserName)
			valueExcel := strconv.FormatFloat(UserTable[i].Contribution2, 'f', -1, 64)
			AddToExcel(WorkSpace+"test1c2.xlsx", indexExcel, valueExcel)
		}
	}
}

//将user得贡献值记录到excel中
func AddToExcel(filepath string, rownum int, newvalue string) {
	var xlFile *xlsx.File
	var sheet2 *xlsx.Sheet
	if PathExists(filepath) == false {
		xlFile = xlsx.NewFile()
		sheet2, _ = xlFile.AddSheet("0")
		_ = xlFile.Save(filepath)
	}
	xlFile, _ = xlsx.OpenFile(filepath)
	sheet2 = xlFile.Sheets[0]
	cell2 := sheet2.Row(rownum).AddCell()
	cell2.Value = newvalue
	_ = xlFile.Save(filepath)
}
