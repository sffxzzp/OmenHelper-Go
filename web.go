package main

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

var VERSION string = "0.8"

type Request struct {
	httpreq *http.Request
	Header  *http.Header
	Client  *http.Client
	Debug   int
	Cookies []*http.Cookie
}
type Response struct {
	R       *http.Response
	content []byte
	text    string
	req     *Request
}
type Header map[string]string
type Params map[string]string
type Datas map[string]string
type Files map[string]string
type Auth []string

func Requests() *Request {
	req := new(Request)
	req.httpreq = &http.Request{
		Method:     "GET",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header = &req.httpreq.Header
	req.Client = &http.Client{}
	jar, _ := cookiejar.New(nil)
	req.Client.Jar = jar
	return req
}
func Get(origurl string, args ...interface{}) (resp *Response, err error) {
	req := Requests()
	resp, err = req.Get(origurl, args...)
	return resp, err
}
func (req *Request) Get(origurl string, args ...interface{}) (resp *Response, err error) {
	req.httpreq.Method = "GET"
	params := []map[string]string{}
	delete(req.httpreq.Header, "Cookie")
	for _, arg := range args {
		switch a := arg.(type) {
		case Header:
			for k, v := range a {
				req.Header.Set(k, v)
			}
		case Params:
			params = append(params, a)
		case Auth:
			req.httpreq.SetBasicAuth(a[0], a[1])
		}
	}
	disturl, _ := buildURLParams(origurl, params...)
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	req.httpreq.URL = URL
	req.ClientSetCookies()
	req.RequestDebug()
	res, err := req.Client.Do(req.httpreq)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp = &Response{}
	resp.R = res
	resp.req = req
	resp.Content()
	defer res.Body.Close()
	resp.ResponseDebug()
	return resp, nil
}
func buildURLParams(userURL string, params ...map[string]string) (string, error) {
	parsedURL, err := url.Parse(userURL)
	if err != nil {
		return "", err
	}
	parsedQuery, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return "", nil
	}
	for _, param := range params {
		for key, value := range param {
			parsedQuery.Add(key, value)
		}
	}
	return addQueryParams(parsedURL, parsedQuery), nil
}
func addQueryParams(parsedURL *url.URL, parsedQuery url.Values) string {
	if len(parsedQuery) > 0 {
		return strings.Join([]string{strings.Replace(parsedURL.String(), "?"+parsedURL.RawQuery, "", -1), parsedQuery.Encode()}, "?")
	}
	return strings.Replace(parsedURL.String(), "?"+parsedURL.RawQuery, "", -1)
}
func (req *Request) RequestDebug() {
	if req.Debug != 1 {
		return
	}
	fmt.Println("===========Go RequestDebug ============")
	message, err := httputil.DumpRequestOut(req.httpreq, false)
	if err != nil {
		return
	}
	fmt.Println(string(message))
	if len(req.Client.Jar.Cookies(req.httpreq.URL)) > 0 {
		fmt.Println("Cookies:")
		for _, cookie := range req.Client.Jar.Cookies(req.httpreq.URL) {
			fmt.Println(cookie)
		}
	}
}
func (req *Request) SetCookie(cookie *http.Cookie) {
	req.Cookies = append(req.Cookies, cookie)
}
func (req *Request) ClearCookies() {
	req.Cookies = req.Cookies[0:0]
}
func (req *Request) ClientSetCookies() {
	if len(req.Cookies) > 0 {
		req.Client.Jar.SetCookies(req.httpreq.URL, req.Cookies)
		req.ClearCookies()
	}
}
func (req *Request) SetTimeout(n time.Duration) {
	req.Client.Timeout = time.Duration(n * time.Second)
}
func (req *Request) Close() {
	req.httpreq.Close = true
}
func (req *Request) Proxy(proxyurl string) {
	urli := url.URL{}
	urlproxy, err := urli.Parse(proxyurl)
	if err != nil {
		fmt.Println("Set proxy failed")
		return
	}
	req.Client.Transport = &http.Transport{
		Proxy:           http.ProxyURL(urlproxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
func (resp *Response) ResponseDebug() {
	if resp.req.Debug != 1 {
		return
	}
	fmt.Println("===========Go ResponseDebug ============")
	message, err := httputil.DumpResponse(resp.R, false)
	if err != nil {
		return
	}
	fmt.Println(string(message))
}
func (resp *Response) Content() []byte {
	var err error
	if len(resp.content) > 0 {
		return resp.content
	}
	var Body = resp.R.Body
	if resp.R.Header.Get("Content-Encoding") == "gzip" && resp.req.Header.Get("Accept-Encoding") != "" {
		reader, err := gzip.NewReader(Body)
		if err != nil {
			return nil
		}
		Body = reader
	}
	resp.content, err = io.ReadAll(Body)
	if err != nil {
		return nil
	}
	return resp.content
}
func (resp *Response) Text() string {
	if resp.content == nil {
		resp.Content()
	}
	resp.text = string(resp.content)
	return resp.text
}
func (resp *Response) SaveFile(filename string) error {
	if resp.content == nil {
		resp.Content()
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(resp.content)
	f.Sync()
	return err
}
func (resp *Response) Json(v interface{}) error {
	if resp.content == nil {
		resp.Content()
	}
	return json.Unmarshal(resp.content, v)
}
func (resp *Response) Cookies() (cookies []*http.Cookie) {
	httpreq := resp.req.httpreq
	client := resp.req.Client
	cookies = client.Jar.Cookies(httpreq.URL)
	return cookies
}
func Post(origurl string, args ...interface{}) (resp *Response, err error) {
	req := Requests()
	resp, err = req.Post(origurl, args...)
	return resp, err
}
func PostJson(origurl string, args ...interface{}) (resp *Response, err error) {
	req := Requests()
	resp, err = req.PostJson(origurl, args...)
	return resp, err
}
func (req *Request) PostJson(origurl string, args ...interface{}) (resp *Response, err error) {
	req.httpreq.Method = "POST"
	req.Header.Set("Content-Type", "application/json")
	delete(req.httpreq.Header, "Cookie")
	for _, arg := range args {
		switch a := arg.(type) {
		case Header:
			for k, v := range a {
				req.Header.Set(k, v)
			}
		case string:
			req.setBodyRawBytes(io.NopCloser(strings.NewReader(arg.(string))))
		case Auth:
			req.httpreq.SetBasicAuth(a[0], a[1])
		default:
			b := new(bytes.Buffer)
			err = json.NewEncoder(b).Encode(a)
			if err != nil {
				return nil, err
			}
			req.setBodyRawBytes(io.NopCloser(b))
		}
	}
	URL, err := url.Parse(origurl)
	if err != nil {
		return nil, err
	}
	req.httpreq.URL = URL
	req.ClientSetCookies()
	req.RequestDebug()
	res, err := req.Client.Do(req.httpreq)
	req.httpreq.Body = nil
	req.httpreq.GetBody = nil
	req.httpreq.ContentLength = 0
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp = &Response{}
	resp.R = res
	resp.req = req
	resp.Content()
	defer res.Body.Close()
	resp.ResponseDebug()
	return resp, nil
}
func (req *Request) Post(origurl string, args ...interface{}) (resp *Response, err error) {
	req.httpreq.Method = "POST"
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	params := []map[string]string{}
	datas := []map[string]string{}
	files := []map[string]string{}
	delete(req.httpreq.Header, "Cookie")
	for _, arg := range args {
		switch a := arg.(type) {
		case Header:
			for k, v := range a {
				req.Header.Set(k, v)
			}
		case Params:
			params = append(params, a)
		case Datas:
			datas = append(datas, a)
		case Files:
			files = append(files, a)
		case Auth:
			req.httpreq.SetBasicAuth(a[0], a[1])
		}
	}
	disturl, _ := buildURLParams(origurl, params...)
	if len(files) > 0 {
		req.buildFilesAndForms(files, datas)
	} else {
		Forms := req.buildForms(datas...)
		req.setBodyBytes(Forms)
	}
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil, err
	}
	req.httpreq.URL = URL
	req.ClientSetCookies()
	req.RequestDebug()
	res, err := req.Client.Do(req.httpreq)
	req.httpreq.Body = nil
	req.httpreq.GetBody = nil
	req.httpreq.ContentLength = 0
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp = &Response{}
	resp.R = res
	resp.req = req
	resp.Content()
	defer res.Body.Close()
	resp.ResponseDebug()
	return resp, nil
}
func (req *Request) setBodyBytes(Forms url.Values) {
	data := Forms.Encode()
	req.httpreq.Body = io.NopCloser(strings.NewReader(data))
	req.httpreq.ContentLength = int64(len(data))
}
func (req *Request) setBodyRawBytes(read io.ReadCloser) {
	req.httpreq.Body = read
}
func (req *Request) buildFilesAndForms(files []map[string]string, datas []map[string]string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for _, file := range files {
		for k, v := range file {
			part, err := w.CreateFormFile(k, v)
			if err != nil {
				fmt.Printf("Upload %s failed!", v)
				panic(err)
			}
			file := openFile(v)
			_, err = io.Copy(part, file)
			if err != nil {
				panic(err)
			}
		}
	}
	for _, data := range datas {
		for k, v := range data {
			w.WriteField(k, v)
		}
	}
	w.Close()
	req.httpreq.Body = io.NopCloser(bytes.NewReader(b.Bytes()))
	req.httpreq.ContentLength = int64(b.Len())
	req.Header.Set("Content-Type", w.FormDataContentType())
}
func (req *Request) buildForms(datas ...map[string]string) (Forms url.Values) {
	Forms = url.Values{}
	for _, data := range datas {
		for key, value := range data {
			Forms.Add(key, value)
		}
	}
	return Forms
}
func openFile(filename string) *os.File {
	r, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return r
}
