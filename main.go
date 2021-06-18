package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"
)

type (
	Config struct {
		Accounts []map[string]string `json:"accounts"`
		LastRun  string              `json:"lastrun"`
	}
)

func run(account map[string]string, timestamp string) bool {
	log.Println("正在运行帐号：" + account["username"])
	lastRun, lastRunStr := getLastRun(timestamp)

	var Login login
	uLogin := Login.newLogin(account["username"], account["password"])
	log.Println("登录准备")
	if !uLogin.webPrepare() {
		return false
	}
	log.Println("开始帐号检查")
	if !uLogin.idpProvider() {
		return false
	}
	log.Println("开始登录")
	localhostUrl := uLogin.webLogin()
	if localhostUrl == "" {
		return false
	}

	var Client client
	uClient := Client.newClient(localhostUrl)
	log.Println("开始模拟 Omen 登录操作")
	if !uClient.clientLogin() {
		return false
	}
	log.Println("正在获取挑战 Session")
	if !uClient.genSession() {
		return false
	}
	log.Println("获取可参与挑战列表")
	challengeList := uClient.getCList()
	if len(challengeList) > 0 {
		log.Println(fmt.Sprintf("正在加入 %d 个任务", len(challengeList)))
		failNum := uClient.joinChallenges(challengeList)
		if failNum == 0 {
			log.Println("任务加入完毕")
		} else if failNum > 0 {
			log.Println(fmt.Sprintf("%d 个任务加入失败", failNum))
		} else {
			log.Fatalln("其他错误")
		}
	}
	log.Println("获取当前待完成任务列表")
	currentList := uClient.getCurCList()
	if currentList == nil {
		log.Fatalln("失败")
		return false
	}
	log.Println(fmt.Sprintf("上次运行时间：%s", lastRunStr))
	currentRun := int(time.Now().Unix())
	if currentRun-lastRun > 3000 && len(challengeList) == 0 {
		log.Println("距上次运行超过 50 分钟，尝试完成所有任务")
		log.Println(fmt.Sprintf("待完成的任务数：%d", len(currentList)))
		uClient.doTask(currentList)
	} else {
		log.Println("距上次运行未超过 50 分钟，尝试完成可立即完成的任务")
		fastList := uClient.getFastList(currentList)
		log.Println(fmt.Sprintf("可立即完成的任务数：%d", len(fastList)))
		uClient.doTask(fastList)
	}
	return true
}

func loadConfig() Config {
	var config Config
	cfgFile := Read("config.json")
	err := json.Unmarshal(cfgFile, &config)
	if err != nil || (err == nil && len(config.Accounts) == 0) {
		log.Fatalln("配置文件读取失败！")
	}
	return config
}

func writeConfig(config Config) bool {
	cfgData, err := json.Marshal(config)
	if err != nil {
		return false
	}
	Write("config.json", string(cfgData))
	return true
}

func exit() {
	fmt.Println()
	log.Println("运行完毕！")
	if runtime.GOOS == "windows" {
		log.Println("请关闭窗口")
		Input("")
	}
}

func main() {
	config := loadConfig()
	for _, account := range config.Accounts {
		if !run(account, config.LastRun) {
			log.Fatalln(fmt.Sprintf("帐号：%s 出现错误", account["username"]))
		}
	}
	config.LastRun = Int2Str(int(time.Now().Unix()))
	if !writeConfig(config) {
		log.Fatalln("写入配置文件时出错！")
	}
	exit()
}
