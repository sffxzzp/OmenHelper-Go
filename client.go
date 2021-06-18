package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type (
	client struct {
		RPCUrl        string
		ApplicationId string
		ClientId      string
		LocalhostUrl  string
		Web           *Request
		AccessToken   string
		TokenType     string
		SessionId     string
	}
	cloginRet struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	tokenRet struct {
		Token string `json:"token"`
	}
	commonRet struct {
		Result map[string]string `json:"result"`
	}
	userInfo struct {
		HpidUserId string `json:"hpid_user_id"`
	}
	challengeListRetCollection struct {
		ChallengeStructureId string            `json:"challengeStructureId"`
		RelevantEvents       []string          `json:"relevantEvents"`
		DisplayName          string            `json:"displayName"`
		Prize                map[string]string `json:"prize"`
		ProgressPercentage   int               `json:"progressPercentage"`
	}
	challengeListRetResult struct {
		Collection []challengeListRetCollection `json:"collection"`
	}
	challengeListRet struct {
		Result challengeListRetResult `json:"result"`
	}
	progressRetResult struct {
		ProgressPercentage int `json:"progressPercentage"`
	}
	progressRet struct {
		Result []progressRetResult `json:"result"`
	}
)

func (c *client) newClient(localhostUrl string) *client {
	Web := Requests()
	Web.Client.Timeout = 30 * time.Second
	return &client{
		RPCUrl:        "https://rpc-prod.versussystems.com/rpc",
		ApplicationId: "6589915c-6aa7-4f1b-9ef5-32fa2220c844",
		ClientId:      "130d43f1-bb22-4a9c-ba48-d5743e84d113",
		LocalhostUrl:  localhostUrl,
		Web:           Web,
	}
}

func (c *client) getFastList(challengeList []map[string]string) []map[string]string {
	fastList := []map[string]string{}
	for _, challenge := range challengeList {
		if challenge["eventName"][:6] == "Launch" {
			fastList = append(fastList, challenge)
		}
	}
	return fastList
}

func (c *client) joinChallenges(challengeList []map[string]string) int {
	num := len(challengeList)
	for _, challenge := range challengeList {
		join := joinData(c.ApplicationId, c.SessionId, challenge["cId"], challenge["csId"])
		res, err := c.Web.PostJson(c.RPCUrl, join.Data)
		if err == nil && res.R.StatusCode < 400 {
			num -= 1
		}
	}
	return num
}

func (c *client) doTask(currentList []map[string]string) {
	for _, challenge := range currentList {
		log.Println(fmt.Sprintf("当前执行：%s - %s%%", challenge["display"], challenge["progress"]))
		time := 45
		if challenge["eventName"][:6] == "Launch" {
			time = 1
		} else {
			time += rand.Intn(20)
		}
		strTime := genTime(time)
		progress := progressData(c.ApplicationId, c.SessionId, strTime, challenge["eventName"])
		res, err := c.Web.PostJson(c.RPCUrl, progress.Data)
		if err != nil {
			continue
		}
		if res.R.StatusCode >= 400 {
			log.Fatalln(fmt.Sprintf("Error: %d", res.R.StatusCode))
			continue
		}
		var progressRet progressRet
		res.Json(&progressRet)
		percentage := progressRet.Result[0].ProgressPercentage
		if percentage == Str2Int(challenge["progress"]) {
			log.Fatalln("进度没有变化，你设置的时间不合理！")
		} else {
			log.Println(fmt.Sprintf("事件：%s|进度：%d%%", challenge["display"], percentage))
		}
	}
}

func (c *client) getCurCList() []map[string]string {
	curList := currentChallengeListData(c.ApplicationId, c.SessionId)
	res, err := c.Web.PostJson(c.RPCUrl, curList.Data)
	if err != nil {
		return nil
	}
	if res.R.StatusCode >= 400 {
		log.Fatalln(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return nil
	}
	var challengeListRet challengeListRet
	res.Json(&challengeListRet)
	challengeList := []map[string]string{}
	for _, item := range challengeListRet.Result.Collection {
		if item.Prize["category"] == "sweepstake" {
			challengeList = append(challengeList, map[string]string{
				"eventName": item.RelevantEvents[0],
				"progress":  Int2Str(item.ProgressPercentage),
				"display":   strings.Split(item.DisplayName, "游戏")[0],
			})
		}
	}
	return challengeList
}

func (c *client) getCList() []map[string]string {
	cList := challengeListData(c.ApplicationId, c.SessionId)
	res, err := c.Web.PostJson(c.RPCUrl, cList.Data)
	if err != nil {
		return nil
	}
	if res.R.StatusCode >= 400 {
		log.Fatalln(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return nil
	}
	var challengeListRet challengeListRet
	res.Json(&challengeListRet)
	challengeList := []map[string]string{}
	for _, item := range challengeListRet.Result.Collection {
		if item.Prize["category"] == "sweepstake" {
			challengeList = append(challengeList, map[string]string{
				"csId":    item.ChallengeStructureId,
				"cId":     item.Prize["campaignId"],
				"event":   item.RelevantEvents[0],
				"display": item.DisplayName,
			})
		}
	}
	return challengeList
}

func (c *client) genSession() bool {
	headers := Header{
		"Authorization": fmt.Sprintf("%s %s", c.TokenType, c.AccessToken),
	}
	tokenUrl := "https://www.hpgamestream.com/api/thirdParty/session/temporaryToken?applicationId=" + c.ApplicationId
	res, err := c.Web.Get(tokenUrl, headers)
	if err != nil {
		return false
	}
	if res.R.StatusCode >= 400 {
		log.Fatalln(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return false
	}
	var tokenRet tokenRet
	res.Json(&tokenRet)
	userToken := tokenRet.Token
	handshake := handshakeData(c.ApplicationId, userToken)
	baseHeaders := Header{
		"Content-Type": "application/json;charset=utf-8",
	}
	res, err = c.Web.PostJson(c.RPCUrl, handshake.Data, baseHeaders)
	if err != nil {
		return false
	}
	if res.R.StatusCode >= 400 {
		log.Fatalln(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return false
	}
	var handshakeRet commonRet
	res.Json(&handshakeRet)
	accountToken := handshakeRet.Result["token"]
	uInfo, err := base64.RawURLEncoding.DecodeString(strings.Split(c.AccessToken, ".")[1])
	if err != nil {
		return false
	}
	var userInfo userInfo
	json.Unmarshal(uInfo, &userInfo)
	start := startData(c.ApplicationId, accountToken, userInfo.HpidUserId)
	res, err = c.Web.PostJson(c.RPCUrl, start.Data, baseHeaders)
	if err != nil {
		return false
	}
	if res.R.StatusCode >= 400 {
		log.Fatalln(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return false
	}
	var startRet commonRet
	res.Json(&startRet)
	c.SessionId = startRet.Result["sessionId"]
	return true
}

func (c *client) clientLogin() bool {
	q, _ := url.Parse(c.LocalhostUrl)
	query := q.Query()
	data := Datas{
		"grant_type":   "authorization_code",
		"code":         query.Get("code"),
		"client_id":    c.ClientId,
		"redirect_uri": q.Scheme + "://" + q.Host + q.Path,
	}
	headers := Header{
		"User-Agent": "Mozilla/4.0 (compatible; MSIE 5.01; Windows NT 5.0)', 'Accept': 'application/json",
		"Except":     "100-continue",
	}
	oauthUrl := "https://oauth.hpbp.io/oauth/v1/token"
	res, err := c.Web.Post(oauthUrl, data, headers)
	if err != nil {
		return false
	}
	if res.R.StatusCode >= 400 {
		log.Fatalln(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return false
	}
	var cloginRet cloginRet
	res.Json(&cloginRet)
	c.AccessToken = cloginRet.AccessToken
	c.TokenType = strings.Title(cloginRet.TokenType)
	return true
}
