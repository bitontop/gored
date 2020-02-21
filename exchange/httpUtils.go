package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func HttpGetRequest(strUrl string, mapParams map[string]string) string {
	httpClient := &http.Client{}

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams := Map2UrlQuery(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}

	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Connection", "close")

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

func HttpPostRequest(strUrl string, mapParams map[string]string) string {
	httpClient := &http.Client{}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "zh-cn")

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

//Signature加密
func ComputeMD5(strMessage string) string {
	h := md5.New()
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum(nil))
}

func ComputeHmacMd5(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(md5.New, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum([]byte("")))
}

func ComputeHmac1(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func ComputeHmac256(strMessage string, strSecret string) string {
	key, _ := base64.StdEncoding.DecodeString(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum(nil))
}

func ComputeBase64Hmac256(strMessage string, strSecret string) string {
	key, _ := base64.StdEncoding.DecodeString(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func ComputeHmac512(strMessage string, strSecret string) string {
	key, _ := base64.StdEncoding.DecodeString(strSecret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum(nil))
}

func ComputeHmac256Base64(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func ComputeHmac256NoDecode(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))
	return hex.EncodeToString(h.Sum(nil))
}

func ComputeHmac256URL(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func ComputeHmac512NoDecode(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum(nil))
}

// 将map格式的请求参数转换为字符串格式的
// mapParams: map格式的参数键值对
// return: 查询字符串
func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += (key + "=" + mapParams[key] + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

func Map2UrlQueryUrl(mapParams map[string]string) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += (key + "=" + url.QueryEscape(mapParams[key]) + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

func Map2UrlQueryInterface(mapParams map[string]interface{}) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += (key + "=" + fmt.Sprintf("%v", mapParams[key]) + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

func GetExternalIP() string {
	httpClient := &http.Client{}

	strRequestUrl := "http://myexternalip.com/raw"

	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}
