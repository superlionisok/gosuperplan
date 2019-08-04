package dbmodels

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
	//	engine.ShowWarn=true
	//需要生成的数据库此处需要加进
	err2 := DB.Sync2(new(LotLottery), new(LotOpenResult))
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

//彩种表
type LotLottery struct {
	//id
	ID int `xorm:"pk autoincr 'ID'"`

	Title      string    `xorm:"'Name'"`
	CreateTime time.Time `xorm:"'CreateTime'"`
	Contents   string    `xorm:"'Contents'"`
	IsDel      int       `xorm:"'IsDel'"`
	AbName     string    `xorm:"'AbName'"`
	//每天开奖次数
	DayOpenCount int `xorm:"'DayOpenCount'"`
	//每天开奖时间
	DayStartTime time.Time `xorm:"'DayStartTime'"`
	//每天结束时间
	DayEndTime time.Time `xorm:"'DayEndTime'"`
	//开奖频率（几分钟一次） 单位 秒
	Rate int `xorm:"'Rate'"`
	//开奖结果延迟时间 单位秒
	Delay         int    `xorm:"'Delay'"`
	OpenResultUrl string `xorm:"'OpenResultUrl'"`
	LotteryTypeID int    `xorm:"'LotteryTypeID'"`
	//期号类型（1每天重复（日期加数字+1），2 无限累加）
	TermNoType int `xorm:"'TimeType'"`
	//每天第一次开奖延迟时间段。和0点做计算 单位 秒
	StartTimeSpan int `xorm:"'StartTimeSpan'"`
}

//开奖结果
type LotOpenResult struct {
	ID         int       `xorm:"pk autoincr 'ID'"`
	LotteryID  int       `xorm:"'LotteryID'"`
	OpenTime   time.Time `xorm:"'OpenTime'"`
	CreateTime time.Time `xorm:"'CreateTime'"`
	TermNumber string    `xorm:"'TermNumber'"`
	OpenNumber string    `xorm:"'OpenNumber'"`
}
