package dbmodels

import (
	"container/list"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/larspensjo/config"
	"log"
	"strconv"
	"time"
)

var DB *xorm.Engine
var (
	configFile = flag.String("configfile", "./conf/conf.ini", "General configuration file")
)
var TOPIC = make(map[string]string)

func getDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:123456@tcp(10.10.15.125:8066)/TESTDB?charset-utf-8&parseTime=true")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db.SetMaxOpenConns(30) //设置最大打开的连接数
	db.SetMaxIdleConns(1)  // 设置最大闲置连接数
	db.Ping()
	return db, nil
}

func Query() {
	customers := list.New()
	db, err := getDB()
	rows, err := db.Query("select * from manager")
	if err != nil {
		log.Println("error-44 : ", err)
	}
	for rows.Next() {
		var ID int
		var Mname string
		var CreateTime time.Time

		if err := rows.Scan(&ID, &Mname, &CreateTime); err != nil {
			log.Println("error-52 : ", err)
		}
		c := Manager{ID, Mname, CreateTime}
		customers.PushBack(c)
	}
	defer rows.Close()
	defer db.Close()
	for cs := customers.Front(); cs != nil; cs = cs.Next() {
		cu := cs.Value.(Manager)
		fmt.Println("custromer : ", cu.ID, cu.Mname, cu.CreateTime)

	}
}
func MyCatInsert() {
	var t0 = time.Now()
	var managers = []Manager{}

	var qidongid = 5850200
	var instertCount = 10000

	for k := 1; k < 101; k++ {
		fmt.Println("第", k, "次插入。qidongid=", qidongid)
		managers = []Manager{}
		for i := 0; i < instertCount; i++ {

			var add = Manager{}
			add.ID = qidongid
			add.Mname = "name"
			add.CreateTime = time.Now()
			managers = append(managers, add)
			qidongid = qidongid + 1
		}
		//insert into manager ( ID,Mname, CreateTime) values(11,'name','2019-08-06' ),(12,'name','2019-08-06' )
		var strSql = "insert into manager (ID,Mname, CreateTime) values"
		var strValue = ""
		var strdate = "2019-01-01"
		for _, v := range managers {
			var str = "(" + strconv.Itoa(v.ID) + ",'" + v.Mname + "','" + strdate + "'),"
			strValue = strValue + str

		}
		var lencount = len(strValue)
		strValue = strValue[0 : lencount-1]

		strSql = strSql + strValue

		var t1 = time.Now()
		db, err := getDB()
		if err != nil {
			fmt.Println("数据库链接错误")
			return
		}
		defer db.Close()
		stm, err := db.Exec(strSql)
		if err != nil {
			log.Println(err)
			return
		}

		var t2 = time.Now()
		var timespan = t2.Sub(t1).Seconds()
		fmt.Println("插入成功,共计"+strconv.Itoa(instertCount)+"条。耗时", timespan, "秒")
		var timespan2 = t2.Sub(t0).Seconds()
		fmt.Println("插入成功,总计耗时", timespan2, "秒")
		fmt.Println(stm)
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

	//fmt.Println(TOPIC)
	//fmt.Println(TOPIC["debug"])

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

type Manager struct {
	ID         int
	Mname      string
	CreateTime time.Time
}
