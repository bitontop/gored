package utils

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

type HttpPost struct {
	URI         string `json:"uri"`
	RequestBody []byte `json:"body"`
	Proxy       string `json:"proxy"`
	Timeout     int    `json:"timeout"`

	//Reference....
	Request  *http.Request
	Response *http.Response
	//........................

	//Output
	ResponseBody []byte
	StatusCode   int
	DebugMode    bool `json:"debug mode"`
	Error        error
}

type HttpGet struct {
	URI     string `json:"uri"`
	Proxy   string `json:"proxy"`
	Timeout int    `json:"timeout"`

	//Reference....
	Request  *http.Request
	Response *http.Response
	//........................

	//Output
	DebugMode    bool `json:"debug mode"`
	ResponseBody []byte
	StatusCode   int
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
	if httpPost.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpPost.Timeout) * time.Second
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
	httpPost.StatusCode = httpPost.Response.StatusCode
	httpPost.ResponseBody, httpPost.Error = ioutil.ReadAll(httpPost.Response.Body)
	if nil != httpPost.Error {
		return httpPost.Error
	}

	return nil
}

func WebOutboundIP() string {
	uris := []string{
		"4.ifcfg.me",
		"alma.ch/myip.cgi",
		"api.infoip.io/ip",
		"api.ipify.org",
		"bot.whatismyipaddress.com",
		"canhazip.com",
		"checkip.amazonaws.com",
		"eth0.me",
		"icanhazip.com",
		"ident.me",
		"ipecho.net/plain",
		"ipinfo.io/ip",
		"ipof.in/txt",
		"ip.tyk.nu",
		"l2.io/ip",
		"smart-ip.net/myip",
		"tnx.nl/ip",
		"wgetip.com",
		"whatismyip.akamai.com",
	}
	num := len(uris)

	for i := 0; i < 10; i++ {
		x := int(rand.Float64() * float64(num))
		strRequestUrl := fmt.Sprintf("http://%s", uris[x])
		log.Printf("Use %s", strRequestUrl)

		get := &HttpGet{
			URI: strRequestUrl,
			// DebugMode: true,
		}
		if err := HttpGetRequest(get); err != nil {
			// log.Printf("ERROR %s", err)
			continue
		}

		ipv4 := string(get.ResponseBody)
		// log.Printf("ipv4 %s", ipv4)
		if net.ParseIP(ipv4) == nil {
			continue
		}
		return ipv4
	}

	return "0.0.0.0"

}

func GetExternalIP() string {

	// httpClient := &http.Client{}

	// strRequestUrl := "http://myexternalip.com/raw"

	// request, err := http.NewRequest("GET", strRequestUrl, nil)
	// if nil != err {
	// 	return "", err
	// }

	// response, err := httpClient.Do(request)
	// if nil != err {
	// 	return "", err
	// }
	// defer response.Body.Close()

	// body, err := ioutil.ReadAll(response.Body)
	// if nil != err {
	// 	return "", err
	// }
	// return string(body), nil
	ip, err := DigOutboundIP()
	if err != nil {
		ip, err := GetOutboundIP()
		if err != nil {
			// return fmt.Sprintf("ERROR: %s", err.Error())
			return WebOutboundIP()
		}
		return fmt.Sprintf("%s", ip)
	}

	return ip
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
	if httpGet.Timeout > 0 {
		httpClient.Timeout = time.Duration(httpGet.Timeout) * time.Second
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
	httpGet.StatusCode = httpGet.Response.StatusCode
	httpGet.ResponseBody, httpGet.Error = ioutil.ReadAll(httpGet.Response.Body)
	if nil != httpGet.Error {
		return httpGet.Error
	}

	return nil
}

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

type OSType string

const (
	Windows OSType = "windows"
	Linux   OSType = "linux"
)

func GetOS() OSType {
	if runtime.GOOS == "windows" {
		return Windows
	} else if runtime.GOOS == "linux" {
		return Linux
	}
	return Linux
}

func DigOutboundIP() (string, error) {
	os := GetOS()

	switch os {
	case Windows:
		// nslookup myip.opendns.com. resolver1.opendns.com
		cmd := "nslookup"
		args := []string{"myip.opendns.com.", "resolver1.opendns.com"}
		out, err := exec.Command(cmd, args...).Output()
		words := strings.Fields(string(out))
		if err == nil {
			return words[len(words)-1], nil
		} else {
			return "", err
		}

	case Linux:

		cmd := "/usr/bin/dig"
		// arg1:="@208.67.222.222"
		// arg2:="ANY"
		// arg3:="myip.opendns.com"
		// arg4:="+short"

		// dig -4 TXT +short o-o.myaddr.l.google.com @ns1.google.com
		// dig -6 TXT +short o-o.myaddr.l.google.com @ns1.google.com //for ipv6
		args := []string{"-4", "TXT", "+short", "o-o.myaddr.l.google.com", "@ns1.google.com"}
		out, err := exec.Command(cmd, args...).Output()
		if err == nil {
			return string(out), nil
		} else {
			return "", err
		}
	}

	return "", fmt.Errorf("Not support system")
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
