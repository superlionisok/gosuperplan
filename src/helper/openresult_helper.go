package helper

import (
	"strconv"
	"strings"
)

//根据开奖结果计算出和值
func CountTotal(res string) int {

	var splits = strings.Split(res, ",")
	var sum = 0
	for _, v := range splits {
		intv, _ := strconv.Atoi(v)
		sum += intv

	}
	return sum
}

//根据开奖结果的和值计算单双
func CountTotalDanShuang(total int) string {

	var res = "单"
	if total%2 == 0 {
		res = "双"
	}
	return res
}
