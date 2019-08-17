package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"helper"
	"log"
	"models"
	"os"
	"os/signal"
	"sync"
)

var kafkaAddress = []string{}

func main() {
	fmt.Println("start")
	getKafkaAddress()
	XiaoFeiKafka()

}
func getKafkaAddress() {
	var address = helper.ConfigGet("kafka", "address")
	kafkaAddress = append(kafkaAddress, address)

}
func XiaoFeiKafka() {
	topic := []string{"lotopenresult"}
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	//广播式消费：消费者1
	go createLong(wg, kafkaAddress, topic, "group-1")
	//广播式消费：消费者2
	//go clusterConsumer2(wg, Address, topic, "group-2")

	wg.Wait()

}

// 支持brokers cluster的消费者
func createLong(wg *sync.WaitGroup, brokers, topics []string, groupId string) {
	fmt.Println("开始一轮生成长龙")

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
				fmt.Fprintf(os.Stdout, "从kafka获取到的数据%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)

				var openresult models.LotOpenResult
				err := json.Unmarshal([]byte(msg.Value), &openresult)
				if err != nil {
					fmt.Println("反序列化开奖结果失败，", err.Error())

				} else {
					//计算长龙
					CountLong(openresult)

					//长龙成功不成功都不影响 kafka 的取消
					//consumer.MarkOffset(msg, "") // mark message as processed
					successes++
				}

			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
	fmt.Println("结束一轮生成长龙")
}

//计算长龙
func CountLong(openresult models.LotOpenResult) bool {

	var lottery models.LotLottery
	has, err := models.DB.Where("ID=?", openresult.LotteryID).Get(&lottery)
	if err != nil {
		fmt.Println("获取彩种失败，", err.Error())
		return false
	}
	if !has {
		fmt.Println("获取彩种没有找到id=", openresult.LotteryID)
		return false
	}
	// 根据彩种类型来生成长龙
	switch {
	//生成快三类的长龙
	case lottery.LotteryTypeID == 1:
		CreateLongByks(openresult)
	}

	return true
}

// 根据快三类来生成长龙
func CreateLongByks(op models.LotOpenResult) {
	var total = helper.CountTotal(op.OpenNumber)
	var danshuang = helper.CountTotalDanShuang(total)
	var listlongs = []models.LotLong{}
	err := models.DB.Where("LotteryID=? ", op.LotteryID).Find(&listlongs)
	if err != nil {
		fmt.Println("从数据库获取长龙信息失败，err=", err.Error())
		return
	}

	//  假如没有长龙信息就新建

	if len(listlongs) < 4 {
		var add1 = models.LotLong{}
		add1.LongCount = 0
		add1.TypeSign = "和值:大"
		add1.Title = "和值:大"
		add1.LotteryID = op.LotteryID
		add1.LotteryTypeID = 1
		add1.Sort = 1
		listlongs = append(listlongs, add1)
		var add2 = models.LotLong{}
		add2.LongCount = 0
		add2.TypeSign = "和值:小"
		add2.Title = "和值:小"
		add2.LotteryID = op.LotteryID
		add2.LotteryTypeID = 1
		add2.Sort = 2
		listlongs = append(listlongs, add2)
		var add3 = models.LotLong{}
		add3.LongCount = 0
		add3.TypeSign = "和值:单"
		add3.Title = "和值:单"
		add3.LotteryID = op.LotteryID
		add3.LotteryTypeID = 1
		add3.Sort = 3
		listlongs = append(listlongs, add3)

		var add4 = models.LotLong{}
		add4.LongCount = 0
		add4.TypeSign = "和值:双"
		add4.Title = "和值:双"
		add4.LotteryID = op.LotteryID
		add4.LotteryTypeID = 1
		add4.Sort = 4
		listlongs = append(listlongs, add4)
		ac, err := models.DB.Insert(&listlongs)
		if err != nil {
			fmt.Println("往数据库插入基础长龙失败，err=", err.Error())
			return
		} else {
			fmt.Println("往数据库插入基础长龙成功ac=，", ac)
		}
		err = models.DB.Where("LotteryID=? ", op.LotteryID).Find(&listlongs)
		if err != nil {
			fmt.Println("从数据库获取长龙信息失败，err=", err.Error())
			return
		}

	}

	for i := 0; i < len(listlongs); i++ {
		if listlongs[i].TypeSign == "和值:大" {
			if total > 10 {
				listlongs[i].LongCount += 1
			} else {
				listlongs[i].LongCount = 0
			}

		} else if listlongs[i].TypeSign == "和值:小" {
			if total > 10 {
				listlongs[i].LongCount = 0
			} else {
				listlongs[i].LongCount += 1
			}
		} else if listlongs[i].TypeSign == "和值:单" {
			if danshuang == "单" {
				listlongs[i].LongCount += 1
			} else {
				listlongs[i].LongCount = 0
			}
		} else if listlongs[i].TypeSign == "和值:双" {
			if danshuang == "单" {
				listlongs[i].LongCount = 0
			} else {
				listlongs[i].LongCount += 1
			}
		}
		// 更新数据库
		affected, err := models.DB.Id(listlongs[i].ID).Cols("LongCount").Update(listlongs[i])
		if err != nil {
			fmt.Println("数据库更新长龙失败")
			return
		}
		fmt.Println("数据库更新长龙成功.返回affected=", affected)
	}

}
