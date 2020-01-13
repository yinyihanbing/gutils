package goutils

import (
	"net/url"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"errors"
	"time"
)

// map参数转换成 a=1&b=2d的格式
func ParamsToStr(params map[string]interface{}) string {
	isFirst := true
	requestUrl := ""
	for k, v := range params {
		if !isFirst {
			requestUrl = requestUrl + "&"
		}
		isFirst = false
		if strings.Contains(k, "_") {
			strings.Replace(k, ".", "_", -1)
		}
		v := fmt.Sprintf("%v", v)
		requestUrl = requestUrl + k + "=" + url.QueryEscape(v)
	}
	return requestUrl
}

// Http POST请求
func HttpPost(requestUrl string, params string, headers map[string]string) (string, error) {
	return HttpSend(requestUrl, "POST", params, headers)
}

// Http POST请求 如果使用 NewRequest 来进行 POST 的表单提交，记得设置头部：req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
func HttpPostByMapParams(requestUrl string, params map[string]interface{}, headers map[string]string) (string, error) {
	return HttpPost(requestUrl, ParamsToStr(params), headers)
}

// Http GET请求
func HttpGet(requestUrl string, headers map[string]string) (string, error) {
	return HttpSend(requestUrl, "GET", "", headers)
}

// Http请求
func HttpSend(requestUrl string, method string, params string, headers map[string]string) (string, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// create a request
	req, err := http.NewRequest(method, requestUrl, strings.NewReader(params))
	if err != nil {
		return "", err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	req.Close = true

	// send JSON to firebase
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	// check status
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("http status error, code=%v", resp.StatusCode))
	}
	defer resp.Body.Close()

	// read data
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
