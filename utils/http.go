package utils

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

type HttpPost struct {
	URI         string `json:"uri"`
	RequestBody []byte `json:"body"`
	Proxy       string `json:"proxy"`

	//Reference....
	Request  *http.Request
	Response *http.Response
	//........................

	//Output
	ResponseBody []byte
	DebugMode    bool `json:"debug mode"`
	Error        error
}

type HttpGet struct {
	URI   string `json:"uri"`
	Proxy string `json:"proxy"`

	//Reference....
	Request  *http.Request
	Response *http.Response
	//........................

	//Output
	DebugMode    bool `json:"debug mode"`
	ResponseBody []byte
	Error        error
}

func HttpPostRequest(httpPost *HttpPost) error {
	httpClient := &http.Client{}

	if httpPost.Proxy != "" {
		proxyUrl, err := url.Parse(httpPost.Proxy)
		if err == nil {
			httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
			if httpPost.DebugMode {
				log.Printf("Apply Proxy @ %s", proxyUrl)
			}
		} else {
			// log.Printf("%+v", err)
			return err
		}
	}

	httpPost.Request, httpPost.Error = http.NewRequest("POST", httpPost.URI, bytes.NewBuffer(httpPost.RequestBody))
	if nil != httpPost.Error {
		return httpPost.Error
	}
	httpPost.Request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	httpPost.Request.Header.Add("Content-Type", "application/json")
	httpPost.Request.Header.Add("Accept-Language", "zh-cn")

	httpPost.Response, httpPost.Error = httpClient.Do(httpPost.Request)
	if nil != httpPost.Error {
		return httpPost.Error
	}
	defer httpPost.Response.Body.Close()

	httpPost.ResponseBody, httpPost.Error = ioutil.ReadAll(httpPost.Response.Body)
	if nil != httpPost.Error {
		return httpPost.Error
	}

	return nil
}

func GetExternalIP() (string, error) {
	httpClient := &http.Client{}

	strRequestUrl := "http://myexternalip.com/raw"

	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return "", err
	}

	response, err := httpClient.Do(request)
	if nil != err {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return "", err
	}

	return string(body), nil
}

func GetInternalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("No Local IP found")
}

func HttpGetRequest(httpGet *HttpGet) error {
	httpClient := &http.Client{}

	if httpGet.Proxy != "" {
		proxyUrl, err := url.Parse(httpGet.Proxy)
		if err == nil {
			httpClient = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
			if httpGet.DebugMode {
				log.Printf("Apply Proxy @ %s", proxyUrl)
			}
		} else {
			// log.Printf("%+v", err)
			return err
		}
	}

	httpGet.Request, httpGet.Error = http.NewRequest("GET", httpGet.URI, nil)
	if nil != httpGet.Error {
		return httpGet.Error
	}
	httpGet.Request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	httpGet.Request.Header.Add("Connection", "close")

	httpGet.Response, httpGet.Error = httpClient.Do(httpGet.Request)
	if nil != httpGet.Error {
		return httpGet.Error
	}
	defer httpGet.Response.Body.Close()

	httpGet.ResponseBody, httpGet.Error = ioutil.ReadAll(httpGet.Response.Body)
	if nil != httpGet.Error {
		return httpGet.Error
	}

	return nil
}
