package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
)

func Input(str string) string {
	var inputStr string
	fmt.Print(str)
	fmt.Scan(&inputStr)
	return inputStr
}

func Read(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func Write(path, data string) {
	err := os.WriteFile(path, []byte(data), 0777)
	if err != nil {
		panic(err)
	}
}

func Int2Str(i int) string {
	return strconv.Itoa(i)
}

func Str2Int(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func loadConfig() []map[string]string {
	var config []map[string]string
	cfgFile := Read("config.json")
	err := json.Unmarshal(cfgFile, &config)
	if err != nil || (err == nil && len(config) == 0) {
		logErr("配置文件读取失败！")
	}
	return config
}

func writeConfig(config []map[string]string) bool {
	cfgData, err := json.Marshal(config)
	if err != nil {
		return false
	}
	Write("config.json", string(cfgData))
	return true
}

func logErr(err string) {
	if err != "" {
		log.Println(err)
	} else {
		log.Println("运行完毕！")
	}
}

func exit() {
	logErr("")
	if runtime.GOOS == "windows" {
		Input("")
	}
}
