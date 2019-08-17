package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"helper"
)

var DB *xorm.Engine

func init() {
	var err error
	var constr = getMysqlConf()
	DB, err = xorm.NewEngine("mysql", constr)
	if err != nil {
		fmt.Println("connet to mysql err,", err.Error())
	}
	DB.ShowSQL(true)
	//需要入库的表就写在下面
	err = DB.Sync2(new(LotOpenResult), new(LotLottery), new(LotLong))

	if err != nil {
		fmt.Println("检查数据库变化出错：", err.Error())
	}

}

func getMysqlConf() string {
	/*
		[mysql]
		name = superplan
		host = 10.10.15.202
		port = 3306
		user = root
		pwd  = 123456

	*/

	var name = helper.ConfigGet("mysql", "name")
	var host = helper.ConfigGet("mysql", "host")
	var port = helper.ConfigGet("mysql", "port")
	var user = helper.ConfigGet("mysql", "user")
	var pwd = helper.ConfigGet("mysql", "pwd")
	// "root:123456@tcp(10.10.15.105:3066)/TESTDB?charset-utf-8&parseTime=true"
	var sqlconstr = user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + name + "?charset-utf-8&parseTime=true"
	return sqlconstr

}
