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

func getLocalUrl(account map[string]string) string {
	var Login login
	uLogin := Login.newLogin(account["username"], account["password"])
	log.Println("登录准备")
	if !uLogin.webPrepare() {
		return ""
	}
	log.Println("开始帐号检查")
	if !uLogin.idpProvider() {
		return ""
	}
	log.Println("开始登录")
	return uLogin.webLogin()
}

func getSessionId(uClient *client, localhostUrl string) bool {
	uClient.setLocalUrl(localhostUrl)
	log.Println("开始模拟 Omen 登录操作")
	if !uClient.clientLogin() {
		return false
	}
	log.Println("正在获取挑战 Session")
	return uClient.genSession()
}

func run(account map[string]string, timestamp string) bool {
	log.Println("正在运行帐号：" + account["username"])
	lastRun, lastRunStr := getLastRun(timestamp)

	var Client client
	uClient := Client.newClient()

	relogin := false
	var clientTest []map[string]string
	var challengeList []map[string]string
	if account["sessionId"] != "" {
		log.Println("检测到帐号 Session，测试可用性")
		uClient.SessionId = account["sessionId"]
		clientTest = uClient.getCList()
		if len(clientTest) > 0 {
			for k, v := range clientTest[0] {
				if k == "Code" && v == "603" {
					log.Fatalln("Session 无效，重新登录")
					localhostUrl := getLocalUrl(account)
					if localhostUrl != "" {
						getSessionId(uClient, localhostUrl)
					}
					relogin = true
				}
			}
		}
	} else {
		localhostUrl := getLocalUrl(account)
		if localhostUrl != "" {
			getSessionId(uClient, localhostUrl)
		}
		relogin = true
	}
	if relogin {
		account["sessionId"] = uClient.SessionId
		log.Println("获取可参与挑战列表")
		challengeList = uClient.getCList()
	} else {
		challengeList = clientTest
	}
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
	if len(currentList) == 0 {
		log.Println("无可完成的任务")
	} else {
		log.Println(fmt.Sprintf("上次运行时间：%s", lastRunStr))
		currentRun := int(time.Now().Unix())
		if currentRun-lastRun > 3000 && len(challengeList) == 0 {
			log.Println("慢速模式，尝试完成所有任务")
			log.Println(fmt.Sprintf("待完成的任务数：%d", len(currentList)))
			uClient.doTask(currentList)
		} else {
			log.Println("快速模式，尝试完成可立即完成的任务")
			fastList := uClient.getFastList(currentList)
			log.Println(fmt.Sprintf("可立即完成的任务数：%d", len(fastList)))
			uClient.doTask(fastList)
		}
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
