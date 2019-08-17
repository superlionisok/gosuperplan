package models

import (
	"time"
)

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

//开奖结果长龙类型
type LotLongType struct {
	ID            int       `xorm:"pk autoincr 'ID'"`
	LotteryTypeID int       `xorm:"'LotteryTypeID'"`
	Title         time.Time `xorm:"'Title'"`
	TypeSign      string    `xorm:"'TypeSign'"`
	Sort          int       `xorm:"'Sort'"`
}

//开奖结果长龙
type LotLong struct {
	ID            int    `xorm:"pk autoincr 'ID'"`
	LotteryTypeID int    `xorm:"'LotteryTypeID'"`
	LotteryID     int    `xorm:"'LotteryID'"`
	LongTypeID    int    `xorm:"'LongTypeID'"`
	Title         string `xorm:"'Title'"`
	TypeSign      string `xorm:"'TypeSign'"`
	Sort          int    `xorm:"'Sort'"`
	LongCount     int    `xorm:"'LongCount'"`
}
