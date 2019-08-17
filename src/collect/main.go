package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"helper"
	"io/ioutil"
	"models"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	MapLotteryTermNo map[int]string
	kafkaAddress     = []string{}
)

func main() {
	fmt.Println("start")
	MapLotteryTermNo = make(map[int]string)
	getKafkaAddress()
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
	fmt.Println(MapLotteryTermNo)
	var list = getAllLottery()
	if list == nil {
		return
	}
	var count = len(list)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for _, v := range list {
		//collectOne(v)

		go collectOneByAPI(&wg, v)
	}

	fmt.Println("采集完成")

	//采集完成。入库
	//affected, err := models.DB.Insert(needInsertOpenResults)
	//if err != nil {
	//	fmt.Println("采集结果入库失败，原因：", err.Error())
	//	return
	//}
	//fmt.Println("采集成功,本次成功入库", affected, "条开奖结果")

}
func collectOneByAPI(wg *sync.WaitGroup, lotteryModel models.LotLottery) {
	defer wg.Done()
	//defer catchPanic("采集" + lotteryModel.AbName + "过程异常")
	//url := "http://10.10.15.202:19001/api/OpenResult?uid=1&lotteryid=" + lotteryModel.Title //请求地址
	url := helper.ConfigGet("url", "apiaddress")

	url = url + strconv.Itoa(lotteryModel.ID)
	var c = http.Client{}
	c.Timeout = time.Second * 10

	resp, err := c.Get(url)
	if err != nil {
		fmt.Println("采集" + lotteryModel.AbName + "过程异常")
		//panic(err)
		return
	}
	//bodyString, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bodyString))
	if resp.StatusCode != 200 {
		fmt.Println("err")
	}
	bodyString, err := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &data)
	if err != nil {
		fmt.Println(err)
	}
	if data["Status"] == "1" {
		datason := data["Data"].([]interface{})

		var datason0 = datason[0].(map[string]interface{})
		var res models.LotOpenResult
		res.CreateTime = time.Now()
		res.OpenNumber = (datason0["OpenNumber"]).(string)
		res.TermNumber = (datason0["TermNumber"]).(string)
		res.LotteryID = lotteryModel.ID
		res.OpenTime = time.Now()

		if MapLotteryTermNo[lotteryModel.ID] == res.TermNumber {
			fmt.Println("记录里已经有该记录，lotterid=", lotteryModel.ID, "TermNumber=", res.TermNumber)
			return
		}

		var instertRes = InsertOneOpenResult(res)
		if instertRes {
			MapLotteryTermNo[lotteryModel.ID] = res.TermNumber
		}
	}

}

func getAllLottery() []models.LotLottery {
	var list []models.LotLottery
	err := models.DB.Where("IsDel=?", 0).Find(&list)
	if err != nil {
		//panic("数据库链接错误")
		fmt.Println("数据库链接错误", time.Now().Format("2006-01-02 15:04:05"))
		return nil
	}
	return list
}

func InsertOneOpenResult(result models.LotOpenResult) bool {

	op := new(models.LotOpenResult)
	has, err := models.DB.Where("LotteryID=? and TermNumber=?", result.LotteryID, result.TermNumber).Get(op)
	if err != nil {
		fmt.Println("查询开奖记录失败,彩种ID=", result.LotteryID, ",TermNumber="+result.TermNumber+",err=", err.Error())
		return false
	}
	if !has {
		id, err := models.DB.Table("lot_open_result").InsertOne(result)

		if err != nil {
			fmt.Println("插入开奖记录失败,彩种ID=", result.LotteryID, ",TermNumber="+result.TermNumber+",err=", err.Error())
			return false
		}
		fmt.Println("插入开奖记录成功,插入ID=", id, ",彩种ID=", result.LotteryID, ",TermNumber="+result.TermNumber)
		result.ID = int(id)
		str, err := json.Marshal(result)
		if err != nil {
			fmt.Println("json话开奖结果失败,彩种ID=", result.LotteryID, ",TermNumber="+result.TermNumber+",err=", err.Error())
		} else {
			ShenChanKafkaYiBu(string(str))
		}

	}
	return true
}

func catchPanic(msg string) {
	err := recover() // recover() 捕获panic异常，获得程序执行权。
	//fmt.Println("发现异常：", msg) // recover()后的内容会正常打印
	if err != nil {
		fmt.Println("发现异常：", msg, err) // runtime error: index out of range
	}
}

func getKafkaAddress() {
	var address = helper.ConfigGet("kafka", "address")
	kafkaAddress = append(kafkaAddress, address)

}
func ShenChanKafkaYiBu(str string) {
	// 新建一个arama配置实例
	config := sarama.NewConfig()

	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true //必须有这个选项
	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Return.Errors = true

	// 新建一个异步生产者

	//	client, err := sarama.NewSyncProducer(Address, config)
	client, err := sarama.NewAsyncProducer(kafkaAddress, config)

	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.AsyncClose()
	//必须有这个匿名函数内容
	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					fmt.Println("err=", err.Error())
				}
			case <-success:
			}
		}
	}(client)

	// 定义一个生产消息，包括Topic、消息内容、
	msg := &sarama.ProducerMessage{}
	msg.Topic = "lotopenresult"
	msg.Key = sarama.StringEncoder("op")
	msg.Value = sarama.StringEncoder(str)
	// 发送消息
	client.Input() <- msg

	fmt.Println("kafka 发送完成", str)
}
