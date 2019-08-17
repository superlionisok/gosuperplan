package dbmodels

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/larspensjo/config"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var DB *xorm.Engine
var (
	configFile = flag.String("configfile", "./conf/conf.ini", "General configuration file")
)
var TOPIC = make(map[string]string)

func init() {

}
func getDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:123456@tcp(10.10.15.105:3306)/superplan?charset=utf-8&parseTime=true")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db.SetMaxOpenConns(30) //设置最大打开的连接数
	db.SetMaxIdleConns(1)  // 设置最大闲置连接数
	db.Ping()
	return db, nil
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
	ID         int `pk autoincr`
	Mname      string
	CreateTime time.Time
}
type Manager2 struct {
	Guid       string
	Mname      string
	Age        int
	GroupName  string
	CreateTime time.Time
}

type dbmodel interface {
	Add(m interface{})
	Del()
}

func Add(m interface{}) {
	/*
	   type:  dbmodels.Manager
	   value:  {1  0001-01-01 00:00:00 +0000 UTC <nil>}
	*/
	fmt.Println("type: ", reflect.TypeOf(m))
	fmt.Println("value: ", reflect.ValueOf(m))
	var r = reflect.TypeOf(m).Name()

	fmt.Println(r)
	//var strsql = " INSERT INTO `lot_long` (`ID`, `LotteryTypeID`, `LotteryID`, `LongTypeID`, `Title`, `TypeSign`, `Sort`, `LongCount`) VALUES ('9', '1', '22', '0', '和值:双', '和值:双', '4', '1') "
	var fieldstr = ""
	var valuestr = ""
	// 拿到model的字段组合
	for i := 0; i < reflect.TypeOf(m).NumField(); i++ {
		fieldName := reflect.TypeOf(m).Field(i).Name
		fieldType := reflect.TypeOf(m).Field(i).Type
		var f = reflect.TypeOf(m).Field(i)
		var tag = reflect.TypeOf(m).Field(i).Tag
		if strings.Index(string(tag), "autoincr") > -1 {
			continue
		}
		fmt.Println(tag, f)
		value := reflect.ValueOf(m).Field(i)

		fieldstr += ", `" + fieldName + "` "
		if fieldType.String() == "time.Time" {
			valuestr += ",'" + fmt.Sprint(value)[0:19] + "'"
		} else {
			valuestr += ",'" + fmt.Sprint(value) + "'"
		}

		fmt.Println("field.Name=", fieldName, "fieldType=", fieldType, "value=", value)
	}

	fmt.Println(fieldstr)
	fmt.Println("valuestr=", valuestr)

	var strsql = " INSERT INTO " + r + " (" + fieldstr[1:] + ") VALUES (" + valuestr[1:] + ") "

	fmt.Println(strsql)
	db, err := sql.Open("mysql", "root:123456@tcp(10.10.15.105:3306)/superplan?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("数据库链接错误，", err.Error())
		panic("数据库链接错误," + err.Error())
	}
	defer db.Close()
	result, err := db.Exec(strsql)
	if err != nil {
		fmt.Println("查询错误，", err.Error())
		panic("查询错误," + err.Error())
	}
	id, err := result.LastInsertId()

	fmt.Println("插入成功。最后插入的ID=", id)

}

func Update(m interface{}) {

	var tableName = reflect.TypeOf(m).Name()

	//var strsql = " INSERT INTO `lot_long` (`ID`, `LotteryTypeID`, `LotteryID`, `LongTypeID`, `Title`, `TypeSign`, `Sort`, `LongCount`) VALUES ('9', '1', '22', '0', '和值:双', '和值:双', '4', '1') "
	var parastr = ""
	var pkstr = ""
	// 拿到model的字段组合
	for i := 0; i < reflect.TypeOf(m).NumField(); i++ {
		fieldName := reflect.TypeOf(m).Field(i).Name
		fieldType := reflect.TypeOf(m).Field(i).Type
		//var f = reflect.TypeOf(m).Field(i)
		var tag = reflect.TypeOf(m).Field(i).Tag
		value := reflect.ValueOf(m).Field(i)
		if strings.Index(string(tag), "pk") > -1 {
			pkstr += "`" + fieldName + "`='" + fmt.Sprint(value) + "'"
		} else {

			parastr += ", `" + fieldName + "`= "
			if fieldType.String() == "time.Time" {
				parastr += "'" + fmt.Sprint(value)[0:19] + "'"
			} else {
				parastr += "'" + fmt.Sprint(value) + "'"
			}

		}
		//fmt.Println("field.Name=", fieldName, "fieldType=", fieldType, "value=", value)
	}

	// UPDATE `Manager` SET `Mname`='zhangsan22', `CreateTime`='2019-08-13 19:06:00' WHERE (`ID`='2') LIMIT 1
	var strsql = " UPDATE " + tableName + " SET " + parastr[1:] + " WHERE (" + pkstr + ") LIMIT 1"

	fmt.Println(strsql)
	db, err := sql.Open("mysql", "root:123456@tcp(10.10.15.105:3306)/superplan?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("数据库链接错误，", err.Error())
		panic("数据库链接错误," + err.Error())
	}
	defer db.Close()
	result, err := db.Exec(strsql)
	if err != nil {
		fmt.Println("查询错误，", err.Error())
		panic("查询错误," + err.Error())
	}
	id, err := result.LastInsertId()

	fmt.Println("修改成功。最后修改的ID=", id)

}

func Del(m interface{}) {

	var tableName = reflect.TypeOf(m).Name()
	var pkstr = ""
	// 拿到model的字段组合
	for i := 0; i < reflect.TypeOf(m).NumField(); i++ {
		fieldName := reflect.TypeOf(m).Field(i).Name
		//fieldType := reflect.TypeOf(m).Field(i).Type

		var tag = reflect.TypeOf(m).Field(i).Tag
		value := reflect.ValueOf(m).Field(i)
		if strings.Index(string(tag), "pk") > -1 {
			pkstr += "`" + fieldName + "`='" + fmt.Sprint(value) + "'"
		}

	}

	// DELETE FROM `Manager` WHERE (`ID`='1')
	var strsql = " DELETE FROM " + tableName + " WHERE (" + pkstr + ") "

	fmt.Println(strsql)
	db, err := sql.Open("mysql", "root:123456@tcp(10.10.15.105:3306)/superplan?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("数据库链接错误，", err.Error())
		panic("数据库链接错误," + err.Error())
	}
	defer db.Close()
	result, err := db.Exec(strsql)
	if err != nil {
		fmt.Println("删除错误，", err.Error())
		panic("删除错误," + err.Error())
	}
	id, err := result.LastInsertId()

	fmt.Println("删除成功。最后修改的ID=", id)

}

func Find(m interface{}) interface{} {

	var tableName = reflect.TypeOf(m).Name()
	var parastr = ""
	var pkstr = ""
	var scanstr = ""
	// 拿到model的字段组合
	for i := 0; i < reflect.TypeOf(m).NumField(); i++ {
		fieldName := reflect.TypeOf(m).Field(i).Name
		//fieldType := reflect.TypeOf(m).Field(i).Type
		//var f = reflect.TypeOf(m).Field(i)
		var tag = reflect.TypeOf(m).Field(i).Tag
		value := reflect.ValueOf(m).Field(i)
		if strings.Index(string(tag), "pk") > -1 {
			pkstr += "`" + fieldName + "`='" + fmt.Sprint(value) + "'"
		} else {
			parastr += ", " + fieldName + ""
			scanstr += ",*" + fieldName
		}

		//fmt.Println("field.Name=", fieldName, "fieldType=", fieldType, "value=", value)
	}

	// DELETE FROM `Manager` WHERE (`ID`='1')
	var strsql = " SELECT * FROM " + tableName + " WHERE (" + pkstr + ") "

	fmt.Println(strsql)
	db, err := sql.Open("mysql", "root:123456@tcp(10.10.15.105:3306)/superplan?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("数据库链接错误，", err.Error())
		panic("数据库链接错误," + err.Error())
	}
	defer db.Close()

	rowsneed, err := db.Query("SELECT * FROM Manager WHERE ID = ?", 3)
	cols, err := rowsneed.Columns()
	if err != nil {
		fmt.Println("rows.Columns错误，", err.Error())
	}
	fmt.Println("rows.Columns()=", cols)

	for i, col := range cols {

		fmt.Println("cols i=", i, ",v=", col)

	}

	interf, err := RowsToIf(*rowsneed)

	if err != nil {
		fmt.Println("RowsToIf错误，", err.Error())
	}
	fmt.Println(interf)

	jsonData, err := json.Marshal(interf)
	if err != nil {
		fmt.Println("json错误，", err.Error())
	}
	fmt.Println("查询后返回的json=", string(jsonData))

	err = json.Unmarshal(jsonData, &m)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println("查询成功。model=", m)
	return m
}

/**时间对象->字符串*/
func Time2Str(t time.Time) string {
	const shortForm = "2006-01-01 15:04:05"

	temp := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	str := temp.Format(shortForm)
	return str
}
func RowsToIf(rowsneed sql.Rows) ([]map[string]interface{}, error) {
	rows := rowsneed

	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength) //临时存储每行数据
	for index, _ := range cache {              //为每一列初始化一个指针
		var a interface{}
		cache[index] = &a
	}
	var list []map[string]interface{} //返回的切片
	for rows.Next() {
		_ = rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			var datatype = reflect.TypeOf(*data.(*interface{})).String()
			fmt.Println(datatype)
			switch {
			case datatype == "int64":
				item[columns[i]] = (*data.(*interface{})).(int64)
				//取实际类型
			case datatype == "time.Time":
				item[columns[i]] = ((*data.(*interface{})).(time.Time)).Format("2006-01-02 15:04:05")
			case datatype == "[]uint8":
				item[columns[i]] = string((*data.(*interface{})).([]uint8))
			default:
				item[columns[i]] = (*data.(*interface{})).(string)

			}

			//item[columns[i]] = *data.(*interface{}) //取实际类型
		}
		list = append(list, item)
	}
	_ = rows.Close()
	return list, nil
}
