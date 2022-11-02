package utils

import (
	"bigrule/common/global"
	"bytes"
	"encoding/json"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

func IsContains(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func IsContainsInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func IsContainsList(items [][]string, item []string) bool {
	for _, eachItem := range items {
		if reflect.DeepEqual(eachItem, item) {
			return true
		}
	}
	return false
}

func AuthGetServiceAddr(serviceName string) (address string) {
	var retryCount int
	for {
		servers, err := global.EtcdReg.GetService(serviceName)
		if err != nil {
			return
		}
		var services []*registry.Service
		for _, value := range servers {
			services = append(services, value)
		}
		next := selector.RoundRobin(services)
		if node, err := next(); err == nil {
			address = node.Address
		}
		if len(address) > 0 {
			return
		}
		retryCount++
		//重试3次为获取返回空
		if retryCount >= 3 {
			return
		}
	}
}

// HttpClient 会话代理
func HttpClient(method string, url string, params []byte, headparams map[string]string) ([]byte, error) {
	timeout := time.Duration(20 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	// set head of request
	for k, v := range headparams {
		request.Header.Set(k, v)
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// PostUrl 统一发送请求
func PostUrl(params interface{}, url string, headext ...map[string]string) (resp []byte, err error) {
	headparams := map[string]string{"Content-Type": "application/json"}
	if len(headext) > 0 {
		for _, hm := range headext {
			for k, v := range hm {
				headparams[k] = v
			}
		}
	}
	resqbyte, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	resp, err = HttpClient("POST", url, resqbyte, headparams)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
