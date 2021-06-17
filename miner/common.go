package miner

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

var httpClient *http.Client
var httpMu sync.Mutex

func GetHttpClient() *http.Client {
	httpMu.Lock()
	defer httpMu.Unlock()
	if httpClient == nil {
		proxy := func(req *http.Request) (*url.URL, error) {
			proxy := GetConfig().Proxy
			if strings.HasSuffix(req.URL.Host, "fofa.so") {
				proxy = GetConfig().GetFofaProxy()
			} else if strings.HasSuffix(req.URL.Host, "github.com") {
				proxy = GetConfig().GetGithubProxy()
			}
			if proxy != "" {
				return url.Parse(proxy)
			}
			return nil, nil
		}
		transport := &http.Transport{Proxy: proxy}
		httpClient = &http.Client{Transport: transport}
	}
	return httpClient
}

func SimpleGet(u string) (body []byte, err error) {
	return CustomGet(u, nil)
}

func CustomGet(u string, preFunc func(r *http.Request)) (body []byte, err error) {
	client := GetHttpClient()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}

	if preFunc != nil {
		preFunc(req)
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	return ioutil.ReadAll(resp.Body)
}

var re *regexp.Regexp

func init() {
	re, _ = regexp.Compile("[a-zA-Z0-9\\-_]+://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")
}

func MatchLinks(content string, prefix string) (links []string) {
	matches := re.FindAllString(content, -1)
	for _, s := range matches {
		if strings.Contains(s, prefix) {
			links = append(links, s)
		}
	}
	return links
}

func RemoveRep(slc []string) []string {
	var result []string
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}
