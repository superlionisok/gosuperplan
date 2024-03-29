package mymodels

import (
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/larspensjo/config"
	"log"
	"time"
)

var DB *xorm.Engine
var (
	configFile = flag.String("configfile", "./conf/conf.ini", "General configuration file")
)
var TOPIC = make(map[string]string)

func init() {
	/*
	   name = lionjihua
	   host = 10.10.15.202
	   port = 3306
	   user = root
	   pwd
	*/
	GetConfig()
	var err error
	var name = TOPIC["name"]
	var host = TOPIC["host"]
	var port = TOPIC["port"]
	var user = TOPIC["user"]
	var pwd = TOPIC["pwd"]

	var sqlconn = user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + name + "?charset=utf8"

	DB, err = xorm.NewEngine("mysql", sqlconn)
	if err != nil {
		fmt.Println("mysql connect err. err=", err.Error())
		return
	}
	DB.ShowSQL(true)
	DB.SetMaxOpenConns(30)
	DB.SetConnMaxLifetime(1)
	//	engine.ShowWarn=true
	err2 := DB.Sync2(new(LotLottery))
	if err2 != nil {
		fmt.Println("mysql sync2 err. err=", err2.Error())
		return
	}

}

func GetConfig() {

	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("mysql") {
		section, err := cfg.SectionOptions("mysql")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mysql", v)
				if err == nil {
					TOPIC[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END

	fmt.Println(TOPIC)
	fmt.Println(TOPIC["debug"])

}

type LotLottery struct {
	ID    int    `xorm:"pk autoincr 'ID'"`
	Title string `xorm:"'Name'"`
	Sort  int    `xorm:"'Sort'"`
}
type LotOpenResult struct {
	ID         int       `xorm:"pk autoincr 'ID'"`
	LotteryID  int       `xorm:"'LotteryID'"`
	OpenTime   time.Time `xorm:"'OpenTime'"`
	CreateTime time.Time `xorm:"'CreateTime'"`
	TermNumber string    `xorm:"'TermNumber'"`
	OpenNumber string    `xorm:"'OpenNumber'"`
}
