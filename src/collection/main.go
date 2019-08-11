package main

import (
	"context"
	"dbmodels"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/go-redis/redis"
	"go.etcd.io/etcd/clientv3"
	"helper"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
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
//var needInsertOpenResults []dbmodels.LotOpenResult

//先阶段内存记录的彩种开奖期号
var MapLotteryTermNo map[string]string

// 完成线程统计
var chanReadyCount chan int

func main() {
	fmt.Println("start")
	defer catchPanic("主线程异常")
	//	TestWaitGroup()
	//SendToKafka()
	//ShenChanKafkaYiBu()
	//SaramaProducer()
	//XiaoFeiKafka()
	DoRedis()
	//ShenChanKafka()
	//doetcd()
	//var n = helper.ConfigGet("mysql", "name")
	//fmt.Println(n)
	//b := helper.ConfigSave("test", "age", "18")
	//fmt.Println(b)

	//InsertMoreTest()
	//UserMycatQuery()
	//dbmodels.MyCatInsert()
	//MapLotteryTermNo = make(map[string]string)
	//timer1 := time.NewTicker(10 * time.Second)
	////Ticker触发
	//for {
	//	select {
	//	case <-timer1.C:
	//		doCollect()
	//	}
	//}
	//fmt.Println("end")
}

func TestWaitGroup() {
	wg := sync.WaitGroup{}
	fmt.Println("wg start")
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go DoSomeThing(i, &wg)
	}
	wg.Wait()
	fmt.Println("wg  end")

}

func DoSomeThing(i int, wg *sync.WaitGroup) {

	fmt.Println("ren wu i=", i)
	wg.Done()
}

func DoRedis() {

	//https://github.com/go-redis/redis
	client := redis.NewClient(&redis.Options{
		Addr:     "10.10.15.105:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	err = client.Set("key2", "2222222222", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exists")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

var Address = []string{"10.10.15.105:9092"}

func SendToKafka() {
	syncProducer(Address)
}
func syncProducer(address []string) {

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()
	topic := "test"
	srcValue := "sync: this is a messagea. index=%d"
	for i := 0; i < 10; i++ {
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(value),
		}
		part, offset, err := p.SendMessage(msg)
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		} else {
			fmt.Fprintf(os.Stdout, value+"发送成功，partition=%d, offset=%d \n", part, offset)
		}
		time.Sleep(2 * time.Second)
	}
}

var wg sync.WaitGroup

func ShenChanKafka() {
	// 新建一个arama配置实例
	config := sarama.NewConfig()

	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll

	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Return.Successes = true

	// 新建一个同步生产者

	client, err := sarama.NewSyncProducer(Address, config)
	//	client, err := sarama.NewAsyncProducer(Address, config)

	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.Close()

	// 定义一个生产消息，包括Topic、消息内容、
	msg := &sarama.ProducerMessage{}
	msg.Topic = "test"
	msg.Key = sarama.StringEncoder("msg1")
	msg.Value = sarama.StringEncoder("msg1 contents 111111111111111111")

	// 发送消息
	pid, offset, err := client.SendMessage(msg)

	msg2 := &sarama.ProducerMessage{}
	msg2.Topic = "test"
	msg2.Key = sarama.StringEncoder("msg2")
	msg2.Value = sarama.StringEncoder("msg1 contents 22222222222222")
	pid2, offset2, err := client.SendMessage(msg2)

	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
	fmt.Printf("pid2:%v offset2:%v\n", pid2, offset2)

}
func ShenChanKafkaYiBu() {
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
	client, err := sarama.NewAsyncProducer(Address, config)

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
	msg.Topic = "test"
	msg.Key = sarama.StringEncoder("msg1yibu")
	msg.Value = sarama.StringEncoder("msg1yibu contents aaaaaaaaaaaaaa")

	// 发送消息

	client.Input() <- msg

	msg2 := &sarama.ProducerMessage{}
	msg2.Topic = "test"
	msg2.Key = sarama.StringEncoder("msg2yibu")
	msg2.Value = sarama.StringEncoder("msg1yibu contents bbbbbbbbbbbbb")

	client.Input() <- msg2

	time.Sleep(time.Second + 5)
	fmt.Println("fa song wan cheng")
}
func XiaoFeiKafka() {
	topic := []string{"test"}
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	//广播式消费：消费者1
	go clusterConsumer(wg, Address, topic, "group-1")
	//广播式消费：消费者2
	go clusterConsumer2(wg, Address, topic, "group-2")

	wg.Wait()

}

// 异步发送kafka
func SaramaProducer() {

	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V0_10_0_1

	fmt.Println("start make producer")
	//使用配置,新建一个异步生产者
	producer, e := sarama.NewAsyncProducer(Address, config)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer producer.AsyncClose()

	//循环判断哪个通道发送过来数据.
	fmt.Println("start goroutine")
	go func(p sarama.AsyncProducer) {
		for {
			select {
			case <-p.Successes():
				//fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
			case fail := <-p.Errors():
				fmt.Println("err: ", fail.Err)
			}
		}
	}(producer)

	var value string
	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		time11 := time.Now()
		value = "this is a  yibu message  " + strconv.Itoa(i) + time11.Format("15:04:05")

		// 发送的消息,主题。
		// 注意：这里的msg必须得是新构建的变量，不然你会发现发送过去的消息内容都是一样的，因为批次发送消息的关系。
		msg := &sarama.ProducerMessage{
			Topic: "test",
		}
		//将字符串转化为字节数组
		msg.Value = sarama.ByteEncoder(value)
		//fmt.Println(value)

		//使用通道发送
		producer.Input() <- msg

		msg2 := &sarama.ProducerMessage{}
		msg2.Topic = "test"
		msg2.Key = sarama.StringEncoder("msg2yibu")
		msg2.Value = sarama.StringEncoder("msg1yibu contents 22222222222222")

		producer.Input() <- msg2

	}
}

// 支持brokers cluster的消费者
func clusterConsumer(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	defer wg.Done()
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}

func clusterConsumer2(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	defer wg.Done()
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				//fmt.Fprintf(os.Stdout, "%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

				fmt.Println("操作2：组名称=", groupId, "，topic=", msg.Topic, "，partition=", msg.Partition, "offset=", msg.Offset, "value=", string(msg.Value))
				consumer.MarkOffset(msg, "") // mark message as processed
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}

const (
	EtcdKey = "/ningxin/backend/logagent/config/192.168.56.1"
)

func doetcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, EtcdKey, "aaaaaa")
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

}

func testtest() {

	req := []byte{1, 2, 3, 4}
	tx8 := make([]byte, 8)
	var tx8len = len(tx8)
	fmt.Println("tx8的len=", tx8len)
	copy(tx8[len(tx8)-len(req):], req)

	fmt.Println(tx8)
}

func catchPanic(msg string) {
	err := recover() // recover() 捕获panic异常，获得程序执行权。
	//fmt.Println("发现异常：", msg) // recover()后的内容会正常打印
	if err != nil {
		fmt.Println("发现异常：", msg, err) // runtime error: index out of range
	}
}

func doCollect() {
	fmt.Println("开始一轮采集", time.Now().Format("15:04:05"))
	chanReadyCount = nil

	var list = getAllLottery()
	if list == nil {
		return
	}
	allLottery = list
	var needCount = len(list)
	chanReadyCount = make(chan int, needCount)
	for _, v := range list {
		//collectOne(v)

		go collectOneByAPI(v)
	}

	for true {
		var ccount = len(chanReadyCount)
		time.Sleep(time.Second)
		//	fmt.Println("chanReadyCount=", &chanReadyCount, "chanReadyCount长度为：", ccount)
		if ccount == needCount {
			break
		}
	}
	fmt.Println("采集完成")

	//采集完成。入库
	//affected, err := dbmodels.DB.Insert(needInsertOpenResults)
	//if err != nil {
	//	fmt.Println("采集结果入库失败，原因：", err.Error())
	//	return
	//}
	//fmt.Println("采集成功,本次成功入库", affected, "条开奖结果")

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
	var lotTitle = lotteryModel.Title
	fmt.Println(lotTitle)

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

	InsertOneOpenResult(res)
	//needInsertOpenResults = append(needInsertOpenResults, res)

	return
}

func collectOneByAPI(lotteryModel dbmodels.LotLottery) {

	defer catchPanic("采集" + lotteryModel.AbName + "过程异常")
	//url := "http://10.10.15.66:19001/api/OpenResult?uid=1&lotteryid=" + lotteryModel.Title //请求地址
	url := helper.ConfigGet("url", "apiaddress")

	url = url + strconv.Itoa(lotteryModel.ID)
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
	bodyString, err := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	err = json.Unmarshal([]byte(bodyString), &data)
	if err != nil {
		fmt.Println(err)
	}
	if data["Status"] == "1" {
		datason := data["Data"].([]interface{})

		var datason0 = datason[0].(map[string]interface{})
		var res dbmodels.LotOpenResult
		res.CreateTime = time.Now()
		res.OpenNumber = (datason0["OpenNumber"]).(string)
		res.TermNumber = (datason0["TermNumber"]).(string)
		res.LotteryID = lotteryModel.ID
		res.OpenTime = time.Now()
		InsertOneOpenResult(res)

	}

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
	has, err := dbmodels.DB.Where("LotteryID=? and TermNumber=?", result.LotteryID, result.TermNumber).Get(op)
	if err != nil {
		fmt.Println("查询开奖记录失败,彩种ID=", result.LotteryID, ",TermNumber="+result.TermNumber+",err=", err.Error())
		return
	}
	if !has {
		_, err := dbmodels.DB.Table("lot_open_result").Insert(result)

		if err != nil {
			fmt.Println("插入开奖记录失败,彩种ID=", result.LotteryID, ",TermNumber="+result.TermNumber+",err=", err.Error())
			return
		}
		fmt.Println("插入开奖记录成功,彩种ID=", result.LotteryID, ",TermNumber="+result.TermNumber)
	}

}

func InsertMoreTest() {

	var c = 10000
	var openrs = []dbmodels.LotOpenResult{}
	for i := 0; i < c; i++ {
		var add dbmodels.LotOpenResult
		add.LotteryID = 1
		add.TermNumber = "123123"
		add.OpenTime = time.Now()
		add.OpenNumber = "1,2,3"
		add.CreateTime = time.Now()
		openrs = append(openrs, add)
	}
	var t1 = time.Now()
	affected, err := dbmodels.DB.Insert(openrs)
	if err != nil {
		fmt.Println("采集结果入库失败，原因：", err.Error())
		return
	}
	var t2 = time.Now()
	var timespan = t2.Sub(t1).Seconds()
	fmt.Println("插入", c, "条数据。共计耗时", timespan, "秒")
	fmt.Println("采集成功,本次成功入库", affected, "条开奖结果")
}

func UserMycatQuery() {

	dbmodels.Query()

}
