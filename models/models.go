package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type MerkleNode struct {
	ID                 int     `orm:"pk;auto"`
	UserName           string  `orm:"column(name)"`
	User2name          string  `orm:"size(300);column(user2name)"`
	ValueOfMerkleNode  float64 `orm:"column(value)"`
	ContentHash        string  `orm:"size(500);column(ch)"`
	TimeStampOfContent string  `orm:"size(500);column(time)"`
	TypeOfThis         int     `orm:"column(type)"` //0表示知识产权、1表示交易
	Blocktag           bool    `orm:"column(tag)"`  //bool类型默认是false
}

type BlockBody struct {
	BlockBody []MerkleNode `orm:"-"`
}

func (p BlockBody) Len() int {
	return len(p.BlockBody)
}

func (p BlockBody) Swap(i, j int) {
	p.BlockBody[i], p.BlockBody[j] = p.BlockBody[j], p.BlockBody[i]
}

type SortByHash struct{ BlockBody }

func (p SortByHash) Less(i, j int) bool {
	return p.BlockBody.BlockBody[i].ContentHash > p.BlockBody.BlockBody[j].ContentHash
}

type Block struct {
	ID            int `orm:"pk;auto"`
	Index         int `orm:"column(uid)"`
	BlockBody     `orm:"-"`
	BlockCreater  string `orm:"size(300);column(creater)"`
	Hash          string `orm:"size(300);column(hash)"`
	PrevBlockHash string `orm:"size(300);column(prehash)"`
	TimeStamp     string `orm:"size(300);column(time)"`
}

type User struct {
	UserName      string  `orm:"column(username);pk"`       //用户名，IP地址
	UserAddTime   string  `orm:"size(300);column(addtime)"` //用户加入网络时间
	UserFileNume  int     `orm:"column(Fnum)"`              //用户知识产权数量
	Contribution  float64 `orm:"size(300);column(con)"`     //Contribution1+Contribution2
	Contribution1 float64 `orm:"size(300);column(con1)"`    //知识产权所得贡献值
	Contribution2 float64 `orm:"size(300);column(con2)"`    //交易所得贡献值
	Contribution3 float64 `orm:"size(300);column(con3)"`    //在线贡献奖励
	Useronlion    int64   `orm:"column(useronlion)"`
	//用户公钥
}

type FileInBlockchain struct {
	ID          int     `orm:"pk;auto"`
	Index       int64   `orm:"column(uid)"`
	Creater     string  `orm:"size(300);column(creater)"`
	Owner       string  `orm:"size(300);column(owner)"`
	FileContent string  `orm:"size(300);column(文件内容)"` //文件内容
	Value       float64 `orm:"column(作品价格)"`           //此处指价格
}

func init() {
	orm.RegisterModel(new(MerkleNode))
	orm.RegisterModel(new(Block))
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(FileInBlockchain))

}
