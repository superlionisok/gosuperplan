package main

import (
	"dbmodels"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"time"

	"net/http"
	_ "strings"
)

/*
流程：
1 从数据库读取所有的彩种，存进map
2 各个彩种获取开奖结果。存进待入库列表
3 入库
*/
//需要采集的彩种列表
var allLottery []dbmodels.LotLottery

// 需要入库的开奖结果
var needInsertOpenResults []dbmodels.LotOpenResult

//先阶段内存记录的彩种开奖期号
var MapLotteryTermNo map[string]string

// 完成线程统计
var chanReadyCount chan int

func main() {
	fmt.Println("start")
	MapLotteryTermNo = make(map[string]string)
	timer1 := time.NewTicker(10 * time.Second)
	//Ticker触发
	for {
		select {
		case <-timer1.C:
			doCollect()
		}
	}
	fmt.Println("end")
}

func doCollect() {
	fmt.Println("开始一轮采集", time.Now().Format("15:04:05"))
	chanReadyCount = nil
	needInsertOpenResults = needInsertOpenResults[0:0]

	var list = getAllLottery()
	if list == nil {
		return
	}
	allLottery = list
	var needCount = len(list)
	chanReadyCount = make(chan int, needCount)
	for _, v := range list {
		go collectOne(v)
	}

	for true {
		var ccount = len(chanReadyCount)
		time.Sleep(time.Second)
		fmt.Println("chanReadyCount=", &chanReadyCount, "chanReadyCount长度为：", ccount)
		if ccount == needCount {
			break
		}
	}

	//采集完成。入库
	affected, err := dbmodels.DB.Insert(needInsertOpenResults)
	if err != nil {
		fmt.Println("采集结果入库失败，原因：", err.Error())
		return
	}
	fmt.Println("采集成功,本次成功入库", affected, "条开奖结果")

}

func AddChanReadyCount() {
	chanReadyCount <- 1
}

func collectOne(lotteryModel dbmodels.LotLottery) {

	//url := "http://trend.gameabchart001.com/gameChart/OG5K3/K3?rowNumType=1&lotteryName=1%E5%88%86%E5%BF%AB3" //请求地址
	url := "http://trend.gameabchart001.com/gameChart/" + lotteryModel.AbName + "/K3?rowNumType=1&lotteryName=" + lotteryModel.Title //请求地址
	var c = http.Client{}
	c.Timeout = time.Second * 10

	resp, err := c.Get(url)
	defer AddChanReadyCount()
	if err != nil {
		panic(err)
	}
	//bodyString, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bodyString))
	if resp.StatusCode != 200 {
		fmt.Println("err")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("goqueryPanic，err=", err.Error())
		panic(err)
	}
	var termNo = ""
	var openNo = ""
	doc.Find(".issue-numbers").Each(func(i int, selection *goquery.Selection) {
		//fmt.Println(selection.Text())
		if i == 29 {
			termNo = selection.Text()
			return
		}
	})
	fmt.Println("MapLotteryTermNo[lotteryModel.AbName]=", MapLotteryTermNo[lotteryModel.AbName])
	if MapLotteryTermNo[lotteryModel.AbName] == termNo {
		return
	} else {
		MapLotteryTermNo[lotteryModel.AbName] = termNo
	}

	doc.Find(".lottery-numbers").Each(func(i int, selection *goquery.Selection) {
		//fmt.Println(selection.Text())
		if i == 29 {
			openNo = selection.Text()
			return
		}
	})

	var res dbmodels.LotOpenResult
	res.CreateTime = time.Now()
	res.OpenNumber = openNo
	res.TermNumber = termNo
	res.LotteryID = lotteryModel.ID
	res.OpenTime = time.Now()
	needInsertOpenResults = append(needInsertOpenResults, res)

	return
}

func getAllLottery() []dbmodels.LotLottery {
	var list []dbmodels.LotLottery
	err := dbmodels.DB.Where("IsDel=?", 0).Find(&list)
	if err != nil {
		//panic("数据库链接错误")
		fmt.Println("数据库链接错误", time.Now().Format("2006-01-02 15:04:05"))
		return nil
	}
	return list
}

func InsertOneOpenResult(result dbmodels.LotOpenResult) {

	op := new(dbmodels.LotOpenResult)
	has, err := dbmodels.DB.Table("LotOpenResult").Where("id=? and TermNumber=?", result.ID, result.TermNumber).Get(op)
	if err != nil {
		fmt.Println("查询开奖记录失败,彩种ID=", result.ID, ",TermNumber="+result.TermNumber+",err=", err.Error())
		return
	}
	if !has {
		id, err := dbmodels.DB.Table("LotOpenResult").Insert(result)

		if err != nil {
			fmt.Println("插入开奖记录失败,彩种ID=", result.ID, ",TermNumber="+result.TermNumber+",err=", err.Error())
			return
		}
		fmt.Println("插入开奖记录成功,彩种ID=", result.ID, ",TermNumber="+result.TermNumber)
	}

}
