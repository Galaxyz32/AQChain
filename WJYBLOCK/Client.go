package main

/**
多节点同步
首先打开server，然后打开client，client连上的时候会发送自身的数据，服务器广播给所有节点，所有节点接收后加入到本地
各个节点输入自己想要保护文件的路径，然后生成MerkleNode广播给全网节点
当收到3个MerkleNode时候，由第一个加入网络的节点产生新区块，并广播全网，各个节点接收信息
下一步工作：
新节点加入时候同步现有区块链
区块链持久化
贡献值计算
真正p2p
GUI

请在我的文档下新建BlockChain文件夹
client与server端不能有传输不接收的情况


init中首先加入输入username，删除username=
将user、file表中相应表示用户的字段换成int
主要改动在handlemessage中，进行自动化处理。
当产生区块后，在D盘会产生名为test_write.xlsx的文件，然后每次新生成区块后会将贡献值记录到该文件
*/
import (
	"AQChain/models"
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"github.com/tealeg/xlsx/xlsx"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var CtbIndex = 1

var NODESNUM = 0
var MNPerBlock int
var WorkSpace = ""
var VerboseFlag *bool
var WebFlag *bool
var WeiXinFlag *bool
var ServerAddrFlag *string
var SilentMode bool
var DBLogPrint bool
var DebugNoDB bool
var IDname int
var TimeType = "2006-01-02 15:04:05"

var MsgWriter *bufio.Writer
var CtbWriter *bufio.Writer
var LogWriter *bufio.Writer
var RecvWriter *bufio.Writer

var OptMutex = &sync.Mutex{}
var FileMutex = &sync.Mutex{}

var Blockchain []models.Block
var username string

var lastBlockOfBlockchain models.Block
var WebType string
var DBPath string
var maxCIndex int
var UserTable []models.User
var FileTable []models.FileInBlockchain
var filemn []string
var row [11]*xlsx.Row
var flag1 bool //判断当前是否正在产生block，这时候所有节点发送mn往后延迟2s，以防止冲突
var flagfilemn bool
var MM models.MerkleNode
var FF models.FileInBlockchain

var UU models.User
var BB models.Block
var OO orm.Ormer

func Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

//判断文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func generateMerkleNode(path string, value float64, username1 string, username2 string, filecontenthash string) models.MerkleNode {
	newMerkleNode := models.MerkleNode{}
	newMerkleNode.UserName = username1
	if path == "" {
		newMerkleNode.ContentHash = filecontenthash
		newMerkleNode.TimeStampOfContent = time.Now().Format(TimeType)
	} else {
		newMerkleNode.TimeStampOfContent = GetFileModTime(path)
		newMerkleNode.ContentHash, _ = SHA256File(path)
	}
	newMerkleNode.ValueOfMerkleNode = value
	newMerkleNode.User2name = username2
	if newMerkleNode.ValueOfMerkleNode == 0 {
		newMerkleNode.TypeOfThis = 0
	} else {
		newMerkleNode.TypeOfThis = 1
	}

	return newMerkleNode
}

//计算文件hash值,经过测试，txt、PNG、PDF、DOCX、MP4、IMG均可
//大小为385MB的视频文件大约3秒计算出结果，效率较高
func SHA256File(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}

	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func SHA256String(str string) (s string) {
	//使用sha256哈希函数
	h := sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	//由于是十六进制表示，因此需要转换
	return hex.EncodeToString(sum)
}

func SHA256Byte(bts []byte) (s string) {
	//使用sha256哈希函数
	h := sha256.New()
	h.Write(bts)
	sum := h.Sum(nil)
	//由于是十六进制表示，因此需要转换
	return hex.EncodeToString(sum)
}

//获取文件修改时间 返回unix时间戳
func GetFileModTime(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Format(TimeType)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Format(TimeType)
	}

	return fi.ModTime().Format(TimeType)
}

func generateNewBlock(b models.BlockBody, blockcreater string) models.Block {
	lastBlockOfBlockchain = Blockchain[len(Blockchain)-1]

	newBlock := models.Block{}
	newBlock.TimeStamp = time.Now().Format(TimeType)
	newBlock.PrevBlockHash = lastBlockOfBlockchain.Hash
	//	newBlock.Hash = "123"
	newBlock.BlockCreater = blockcreater
	newBlock.Index = lastBlockOfBlockchain.Index + 1
	newBlock.BlockBody = b

	msg, _ := json.Marshal(newBlock)
	newBlock.Hash = SHA256Byte(msg)
	return newBlock
}

func gerateGensisBlocks() models.Block {

	gensisBlock := models.Block{}
	//	gensisBlock.Hash = "123456"
	gensisBlock.Index = 0
	gensisBlock.PrevBlockHash = ""
	gensisBlock.TimeStamp = time.Now().Format(TimeType)

	a := generateMerkleNode(WorkSpace+"/tmp/使用指南.txt", 0, "admin", "", "")
	fileg := models.FileInBlockchain{}
	fileg.FileContent = "111111"
	fileg.Creater = "admin"
	fileg.Value = 0
	fileg.Owner = ""
	FileTable = append(FileTable, fileg)
	blockbody := models.BlockBody{}
	blockbody.BlockBody = append(blockbody.BlockBody, a)
	gensisBlock.BlockBody = blockbody

	//msg, _ := json.Marshal(gensisBlock)
	gensisBlock.Hash = "GENSISBLOCKHASH"

	Blockchain = append(Blockchain, gensisBlock)

	return gensisBlock
}

//清空BlockBody中MerkleNode
func SliceClear2(s *[]models.MerkleNode) {
	*s = (*s)[0:0]
}

//清空BlockBody中MerkleNode
func SliceClear3(f *[]string) {
	*f = (*f)[0:0]
}

//生成区间随机数，randint是从min到max-1取值
func RandInt(min int, max int) int {
	rand.Seed(time.Now().Unix())
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func initDB() {
	var space1 string
	var space2 string
	space1 = CaculatePath()
	space2 = space1 + strconv.Itoa(IDname) + ".db"
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", space2)

	//orm.RegisterDataBase("default","sqlite3","C:/Users/Administrator/go/src/DBTest1/data/test2.db")
	orm.RunSyncdb("default", false, true)
}

	func main() {
	VerboseFlag = flag.Bool("v", false, "啰嗦模式,true则全显示")
	WebFlag = flag.Bool("h", false, "网页模式")
	WeiXinFlag = flag.Bool("wx", false, "微信模式")
	ServerAddrFlag = flag.String("ip", "", "连线服务器地址")

	flag.Parse()
	DBLogPrint = true  //是否显示数据库有关的输出
	DebugNoDB = false  //调试阶段暂不加入数据库
	SilentMode = false //true则不显示任何输出

	WorkSpace = CaculatePath()

	os.Mkdir(WorkSpace+"/tmp", os.ModePerm)
	os.OpenFile(WorkSpace+"/tmp/使用指南.txt", os.O_WRONLY|os.O_CREATE, 0666)

	//RSAKeyGenInit()

	//fmt.Println("我的公钥是")
	//spew.Dump(string(PubKey))

	PubKeyMap = make(map[string]string)

	P2PConnect()
	go updateNodeListFromServer()

	initDB()

	//RSAKeyMap()

	//spew.Dump(PubKeyMap)

	if *WeiXinFlag == true {
		go ListenWXServer()
	}

	_ = godotenv.Load("example.env")
	NODESNUM = getIntFromEnv("NodesNum")
	MNPerBlock = getIntFromEnv("mnPerBlock")

	RecvFile, outputError := os.OpenFile(WorkSpace+"Recv "+strconv.Itoa(IDname)+".log", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with LogFile opening or creation\n")
		return
	}
	//defer outputFile.Close()
	RecvWriter = bufio.NewWriter(RecvFile)

	LogFile, outputError := os.OpenFile(WorkSpace+"Log "+strconv.Itoa(IDname)+".log", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with LogFile opening or creation\n")
		return
	}
	//defer outputFile.Close()
	LogWriter = bufio.NewWriter(LogFile)

	//记录接收到的所有消息
	MsgFile, outputError := os.OpenFile(WorkSpace+"Msg "+strconv.Itoa(IDname)+".dat", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with msgfile opening or creation\n")
		return
	}
	//defer outputFile.Close()
	MsgWriter = bufio.NewWriter(MsgFile)

	//记录贡献值情况
	CtbFile, outputError := os.OpenFile(WorkSpace+"Ctb "+strconv.Itoa(IDname)+".dat", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with ctbfile opening or creation\n")
		return
	}
	//defer outputFile.Close()

	CtbWriter = bufio.NewWriter(CtbFile)

	time.Sleep(1e9)
	WorkSpace = CaculatePath()

	OO = orm.NewOrm() //实例化OO
	OO.Using("default")
	orm.Debug = true //启动数据库orm

	//go DBSelect(MM, OO)   //merklenode表查询
	//go Fileselect(FF, OO) //文件表查询
	ReadBlockFromDB(BB, OO)

	fmt.Printf("\x1b[33m")
	spew.Dump(Blockchain)
	fmt.Printf("\x1b[0m")

	if len(Blockchain) == 0 {
		gerateGensisBlocks()
	}

	//go func() {
	//	for {
	//		spew.Dump(MyFileList)
	//		time.Sleep(3e9)
	//
	//	}
	//}()

	//time.Sleep(10e9)
	OpenTheWeb(MM, OO, UU, BB, FF)
	beego.Run()

}
