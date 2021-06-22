package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	login struct {
		Email         string
		Password      string
		ApplicationId string
		ClientID      string
		IdpProvider   string
		Web           *Request
		Headers       map[string]string
		EntryUrl      string
		BackendCsrf   string
	}
	bRet struct {
		EntryUrl    string `json:"regionEndpointUrl"`
		BackendCsrf string `json:"csrfToken"`
	}
	idpRet struct {
		Error      string              `json:"error"`
		Identities []map[string]string `json:"identities"`
	}
	wlRet struct {
		Status  string `json:"status"`
		NextUrl string `json:"nextUrl"`
	}
)

func (l *login) newLogin(email string, password string) *login {
	Web := Requests()
	Web.Client.Timeout = 30 * time.Second
	return &login{
		Email:         email,
		Password:      password,
		ApplicationId: "6589915c-6aa7-4f1b-9ef5-32fa2220c844",
		ClientID:      "130d43f1-bb22-4a9c-ba48-d5743e84d113",
		IdpProvider:   "hpid",
		Web:           Web,
		Headers: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
			"User-Agent":   "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36",
		},
		EntryUrl:    "https://ui-backend.us-west-2.id.hp.com/bff/v1",
		BackendCsrf: "",
	}
}

func (l *login) webLogin() string {
	url := l.EntryUrl + "/session/username-password"
	res, err := l.Web.PostJson(url, Datas{"username": fmt.Sprintf("%s@%s", l.Email, l.IdpProvider), "password": l.Password})
	if err != nil {
		return ""
	}
	if res.R.StatusCode >= 400 {
		log.Println(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return ""
	}
	var wlRet wlRet
	res.Json(&wlRet)
	if wlRet.Status != "success" {
		log.Println("登录失败")
		return ""
	}
	checkRedir := l.Web.Client.CheckRedirect
	l.Web.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	res1, err := l.Web.Get(wlRet.NextUrl)
	l.Web.Client.CheckRedirect = checkRedir
	if err != nil {
		log.Println(err)
		return ""
	}
	localhostUrl, err := res1.R.Location()
	if err != nil {
		log.Println(err)
		return ""
	}
	return localhostUrl.String()
}

func (l *login) idpProvider() bool {
	url := l.EntryUrl + "/session/check-username"
	l.Headers["csrf-token"] = l.BackendCsrf
	for k, v := range l.Headers {
		l.Web.Header.Set(k, v)
	}
	res, err := l.Web.PostJson(url, Datas{"username": l.Email})
	if err != nil {
		return false
	}
	if res.R.StatusCode >= 400 {
		log.Println(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return false
	}
	var idpRet idpRet
	res.Json(&idpRet)
	if idpRet.Error == "captchaRequired" {
		log.Println("请求过于频繁，需要验证码")
		log.Println("请访问 https://myaccount.id.hp.com/uaa 并登录一次后再试")
	} else {
		if len(idpRet.Identities) == 0 {
			log.Println("检查帐号出错！")
		} else if idpRet.Identities[0]["idpProvider"] != "" {
			log.Println("成功")
			l.IdpProvider = idpRet.Identities[0]["idpProvider"]
			log.Println(fmt.Sprintf("ID 类型：%s", l.IdpProvider))
			log.Println(fmt.Sprintf("区域：%s", idpRet.Identities[0]["locale"]))
			return true
		}
	}
	return false
}

func (l *login) webPrepare() bool {
	url := "https://oauth.hpbp.io/oauth/v1/auth?response_type=code&client_id=" + l.ClientID + "&redirect_uri=http://localhost:9081/login&scope=email+profile+offline_access+openid+user.profile.write+user.profile.username+user.profile.read&state=G5g495-R4cEE" + strconv.FormatFloat(rand.Float64()*100000, 'f', 11, 64) + "&max_age=28800&acr_values=urn:hpbp:hpid&prompt=consent"
	res, err := l.Web.Get(url)
	if err != nil {
		return false
	}
	if res.R.StatusCode >= 400 {
		log.Println(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return false
	}
	backendUrl := "https://ui-backend.id.hp.com/bff/v1/auth/session"
	location := strings.Split(res.R.Request.URL.String(), "=")[1]
	for k, v := range l.Headers {
		l.Web.Header.Set(k, v)
	}
	res, err = l.Web.PostJson(backendUrl, Datas{"flow": location})
	if err != nil {
		return false
	}
	if res.R.StatusCode >= 400 {
		log.Println(fmt.Sprintf("Error: %d", res.R.StatusCode))
		return false
	}
	var bRet bRet
	res.Json(&bRet)
	l.EntryUrl = bRet.EntryUrl
	l.BackendCsrf = bRet.BackendCsrf
	return true
}
