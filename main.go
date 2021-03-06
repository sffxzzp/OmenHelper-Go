package main

import (
	"fmt"
	"log"
	"time"
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

func run(account map[string]string) bool {
	log.Println("正在运行帐号：" + account["username"])

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
					log.Println("Session 无效，重新登录")
					localhostUrl := getLocalUrl(account)
					if localhostUrl != "" {
						getSessionId(uClient, localhostUrl)
					} else {
						return false
					}
					relogin = true
					break
				}
			}
		}
	} else {
		localhostUrl := getLocalUrl(account)
		if localhostUrl != "" {
			getSessionId(uClient, localhostUrl)
		} else {
			return false
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
		log.Printf("正在加入 %d 个任务\n", len(challengeList))
		failNum := uClient.joinChallenges(challengeList)
		if failNum == 0 {
			log.Println("任务加入完毕")
		} else if failNum > 0 {
			log.Printf("%d 个任务加入失败\n", failNum)
		} else {
			log.Println("其他错误")
		}
	}
	log.Println("获取当前待完成任务列表")
	currentList := uClient.getCurCList()
	if currentList == nil {
		log.Println("失败")
		return false
	}
	if len(currentList) == 0 {
		log.Println("无可完成的任务")
	} else {
		if len(challengeList) == 0 {
			log.Println("慢速模式，尝试完成所有任务")
			log.Printf("待完成的任务数：%d\n", len(currentList))
			uClient.doTask(currentList)
		} else {
			log.Println("快速模式，尝试完成可立即完成的任务")
			fastList := uClient.getFastList(currentList)
			log.Printf("可立即完成的任务数：%d\n", len(fastList))
			uClient.doTask(fastList)
		}
	}
	return true
}

func main() {
	accounts := loadConfig()
	for _, account := range accounts {
		retries := 3
		for retries > 0 {
			if run(account) {
				time.Sleep(10 * time.Second)
				break
			}
			log.Printf("帐号：%s 出现错误，即将重试\n\n", account["username"])
			retries -= 1
			time.Sleep(3 * time.Second)
		}
	}
	if !writeConfig(accounts) {
		log.Println("写入配置文件时出错！")
	}
	fmt.Println()
	exit()
}
