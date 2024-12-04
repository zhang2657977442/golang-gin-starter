package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-gin-starter/entity"
	"github.com/golang-gin-starter/utils/log"
)

const (
	POST_METHOD   = "POST"
	GET_METHOD    = "GET"
	DELETE_METHOD = "DELETE"
)

var (
	httpClient = &http.Client{
		Timeout: 60 * time.Second,
	}

	httpClientShort = &http.Client{
		Timeout: time.Second,
	}
)

func SendHttpPostProto(url string, reqBytes []byte) ([]byte, error) {
	reqReader := bytes.NewReader(reqBytes)
	req, err := http.NewRequest(POST_METHOD, url, reqReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/proto")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return respBytes, nil
}

func SendHttpGet(url string, args url.Values) ([]byte, error) {
	var query []byte
	for k, v := range args {
		var t string
		for _, value := range v {
			t = value
		}
		query = append(query, []byte(fmt.Sprintf("%v=%v&", k, t))...)
	}
	req, err := http.NewRequest("GET", url+"?"+string(query), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/proto")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return respBytes, nil
}

func SendGetReqForJson(host, uri string, headerArgs, args url.Values, result interface{}) error {
	return SendReqForJson(host, uri, GET_METHOD, headerArgs, args, result)
}

func SendPostReqForJson(host, uri string, args url.Values, result interface{}) error {

	headerArgs := make(url.Values, 1)
	headerArgs.Set("Content-type", "application/x-www-form-urlencoded")

	return SendReqForJson(host, uri, POST_METHOD, headerArgs, args, result)
}

func SendPostJsonReq(host, uri string, headerArgs url.Values, param interface{}, result interface{}) error {
	if headerArgs == nil {
		headerArgs = make(url.Values)
	}
	headerArgs.Set("Content-type", "application/json")

	reqData, err := json.Marshal(param)
	if err != nil {
		return err
	}

	var (
		url      []string
		req      *http.Request
		finalUrl string
	)

	url = append(url, host)
	if !strings.HasPrefix(uri, "/") {
		url = append(url, "/")
	}
	url = append(url, uri)
	finalUrl = strings.Join(url, "")

	req, err = http.NewRequest(POST_METHOD, finalUrl, bytes.NewReader(reqData))
	if err != nil {
		log.Error("Http Post Send", "Http", "", "post %v with param %v create req failed, err:[%v]", finalUrl, reqData, err)
		return entity.HTTP_REQUEST_ERROR
	}

	if len(headerArgs) != 0 {
		for k, v := range headerArgs {
			req.Header.Set(k, fmt.Sprintf("%v", v[0]))
		}
	}
	client := httpClient
	if strings.Contains(uri, "visual") {
		client = httpClientShort
	}
	tStart := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(tStart).Seconds()

	if err != nil {
		log.Error("Http Post Send", "Http", "", "post %v failed, elapsed %v seconds, err:[%v]", finalUrl, elapsed, err)
		return entity.HTTP_REQUEST_ERROR
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Http Post Send", "Http", "", "post %v response status code[%v], elapsed %v seconds.", finalUrl, resp.StatusCode, elapsed)
		return entity.HTTP_REQUEST_ERROR
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error("Http Post Send", "Http", "", "post %v read body error, elapsed %v seconds, err:[%v]", finalUrl, elapsed, err)
		return entity.HTTP_REQUEST_ERROR
	}

	if err := json.Unmarshal(data, result); err != nil {
		log.Error("Http Post Send", "Http", "", "post %v:%v cannot parse response body[%v] to json.", host, uri, string(data))
		return entity.INTERNAL_ERROR
	}
	return nil
}

func SendJdosReq(url, req_type, erp, token string, param interface{}, result interface{}) error {
	jsonStr, err := json.Marshal(param)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(req_type, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("erp", erp)
	req.Header.Add("token", token)

	tStart := time.Now()
	resp, err := httpClient.Do(req)
	elapsed := time.Since(tStart).Seconds()

	if err != nil {
		log.Error("Http Post Send", "Http", "", "httpPost %v failed, elapsed %v seconds, err:[%v]", url, elapsed, err)
		return entity.HTTP_REQUEST_ERROR
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Http Post Send", "Http", "", "httpPost %v response status code[%v], elapsed %v seconds.", url, elapsed, resp.StatusCode)
		return entity.HTTP_REQUEST_ERROR
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error("Http Post Send", "Http", "", "httpGet %v response read body error, elapsed %v seconds, err:[%v]", url, elapsed, err)
		return entity.HTTP_REQUEST_ERROR
	}

	log.Error("Http Post Send", "Http", "", "httpGet %v response body:[%v], elapsed %v seconds.", url, string(body), elapsed)
	if err := json.Unmarshal(body, result); err != nil {
		log.Error("Http Post Send", "Http", "", "httpGet %v cannot parse response body[%v] to json.", url, string(body))
		return entity.INTERNAL_ERROR
	}

	return nil
}

func SendJdosGetReq(url, erp, token string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("erp", erp)
	req.Header.Add("token", token)

	tStart := time.Now()
	resp, err := httpClient.Do(req)
	elapsed := time.Since(tStart).Seconds()

	if err != nil {
		log.Error("Http Post Send", "Http", "", "httpGet %v failed, elapsed %v seconds, err:[%v]", url, elapsed, err)
		return entity.HTTP_REQUEST_ERROR
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Http Post Send", "Http", "", "httpGet %v response status code[%v], elapsed %v seconds.", url, elapsed, resp.StatusCode)
		return entity.HTTP_REQUEST_ERROR
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error("Http Post Send", "Http", "", "httpGet %v response status code[%v], elapsed %v seconds.", url, elapsed, resp.StatusCode)
		return entity.HTTP_REQUEST_ERROR
	}

	log.Debug("Http Post Send", "Http", "", "httpGet %v response body:[%v], elapsed %v seconds.", url, string(body), elapsed)
	if err := json.Unmarshal(body, result); err != nil {
		log.Error("Http Post Send", "Http", "", "httpGet %v cannot parse response body[%v] to json.", url, string(body))
		return entity.INTERNAL_ERROR
	}

	return nil
}

func SendReqForJson(host, uri, method string, headerArgs, bodyArgs url.Values, result interface{}) error {
	data, err := SendReq(host, uri, method, headerArgs, bodyArgs)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, result); err != nil {
		log.Error("Http Post Send", "Http", "", "%v %v:%v cannot parse response body[%v] to json.", method, host, uri, string(data))
		return entity.INTERNAL_ERROR
	}

	return nil
}

func SendReq(host, uri, method string, headerArgs, bodyArgs url.Values) ([]byte, error) {
	var (
		url      []string
		req      *http.Request
		err      error
		finalUrl string
		//bodyReader *bytes.Reader
	)
	url = append(url, host)
	if len(uri) > 0 && !strings.HasPrefix(uri, "/") {
		url = append(url, "/")
	}
	url = append(url, uri)

	switch method {
	case "POST":
		finalUrl = strings.Join(url, "")
		if len(bodyArgs) != 0 {
			req, err = http.NewRequest(method, finalUrl, bytes.NewReader([]byte(bodyArgs.Encode())))
		} else {
			req, err = http.NewRequest(method, finalUrl, nil)
		}

	case "GET":
		if len(bodyArgs) != 0 {
			url = append(url, "?")
			//for k, v := range bodyArgs {
			//	url = append(url, fmt.Sprintf("&%s=%v&", k, v))
			//}
			url = append(url, bodyArgs.Encode())
		}
		finalUrl = strings.Join(url, "")
		req, err = http.NewRequest(method, finalUrl, nil)

	case "DELETE":
		if len(bodyArgs) != 0 {
			url = append(url, "?")
			url = append(url, bodyArgs.Encode())
		}
		finalUrl = strings.Join(url, "")
		req, err = http.NewRequest(method, finalUrl, nil)
	}

	if err != nil {
		log.Error("Http Post Send", "Http", "", "%v %v with param %v create req failed, err:[%v]", method, finalUrl, bodyArgs, err)
		return nil, entity.HTTP_REQUEST_ERROR
	}

	if len(headerArgs) != 0 {
		for k, v := range headerArgs {
			req.Header.Set(k, fmt.Sprintf("%v", v[0]))
		}
	}

	client := httpClient
	if strings.Contains(uri, "visual") {
		client = httpClientShort
	}
	tStart := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(tStart).Seconds()

	if err != nil {
		log.Error("Http Post Send", "Http", "", "%v %v failed, elapsed %v seconds, err:[%v]", method, finalUrl, elapsed, err)
		return nil, entity.HTTP_REQUEST_ERROR
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("Http Post Send", "Http", "", "%v %v response status code[%v], elapsed %v seconds.", method, finalUrl, resp.StatusCode, elapsed)
		return nil, entity.HTTP_REQUEST_ERROR
	}

	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Error("Http Post Send", "Http", "", "%v %v read body error, elapsed %v seconds, err:[%v]", method, finalUrl, elapsed, err)
		return nil, entity.HTTP_REQUEST_ERROR
	}
	return data, nil
}
