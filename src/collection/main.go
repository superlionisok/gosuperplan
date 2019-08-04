package main

import(
	"dbmodels"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"net/http"
	"strings"
)


func main() {
	fmt.Println("start")



	collectOne()



	fmt.Println("end")
}
func collectOne()  dbmodels.LotOpenResult{

	url := "https://www.600wcp.com/ssc/ajaxGetHistory.json?timestamp=1564745106126" //请求地址





var res dbmodels.LotOpenResult

return  res
}
func httpDo() {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://www.fcyl7.com/tools/_ajax/getLotteryOpenNewestGame", strings.NewReader("name=cjb"))
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("页面打开错误：",err.Error())
	}

	fmt.Println(string(body))
}


func DoPost() {

	url := "http://trend.gameabchart001.com/gameChart/OG5K3/K3?rowNumType=1&lotteryName=1%E5%88%86%E5%BF%AB3" //请求地址
	//contentType := "application/x-www-form-urlencoded; charset=UTF-8"
	//参数，多个用&隔开

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("页面内容",string(body))

	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		fmt.Println("序列化失败", err)
	}
	//{
	//    "result": 1,
	//    "sscHistoryList": [
	//        {
	//            "playGroupName": "分分时时彩",
	//            "playGroupId": 15,
	//            "number": "201908010756",
	//            "openCode": "9,8,2,4,2",
	//            "openTime": 1564634100000,
	//            "date": "2019-08-01"
	//        }
	//    ]
	//}
	m := f.(map[string]interface{})

	var r1 = m["result"]
	fmt.Println(r1)
	var r2 = m["sscHistoryList"]
	fmt.Println(r2)
	var datas = r2.([]interface{})
	var d = datas[0].(map[string]interface{})
	var number = d["number"]
	var opencode = d["openCode"]
	var opentime = d["openTime"]

	var a = opentime.(float64)

	fmt.Println(a)

	var aaa = int64(a)
	//string =string[0:10]
	fmt.Println(aaa)
	bbb := aaa / 1000

	var t = time.Unix(bbb, 0)
	fmt.Println(number, opencode, t)
}
