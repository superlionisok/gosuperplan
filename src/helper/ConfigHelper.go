package helper

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

func ConfigGet(section, key string) string {

	cfg, err := ini.Load("./conf/conf.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return ""
	}

	// 典型读取操作，默认分区可以使用空字符串表示
	//fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
	//fmt.Println("Data Path:", cfg.Section(section).Key(key).String())
	//
	return cfg.Section(section).Key(key).String()

}

func ConfigSave(section, key, value string) bool {

	cfg, err := ini.Load("./conf/conf.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
		return false
	}
	// 修改某个值然后进行保存
	cfg.Section(section).Key(key).SetValue(value)
	cfg.SaveTo("./conf/conf.ini")
	return true
}
