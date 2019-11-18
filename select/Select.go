package main

import (
	"AQChain/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

func SelectMN() {
	var i int = 1
	for i < 1000 {
		o := orm.NewOrm()
		merkle := models.MerkleNode{}
		merkle.UserName = "1"
		//4.查询
		err := o.Read(&merkle)
		if err != nil {
			beego.Info("查询失败查询失败查询失败查询失败查询失败查询失败查询失败查询失败查询失败查询失败失败查询失败查失败查询失败查", err)
			return
		}
		beego.Info("查询成功查询成功查询成功查询成功查询成功查询成功查询成功查询成功查询成功查询成功查询成功查询成功", merkle)
		time.Sleep(10e9)
		i++
	}

}

func main() {
	SelectMN()
}
