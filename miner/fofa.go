package miner

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hluwa/simplethreadpool"
	"html"
	"log"
	"sync"
)

type FofaResponse struct {
	Error  bool   `json:"error"`
	Errmsg string `json:"errmsg"`
}

type FofaSearchResult struct {
	FofaResponse
	Mode    string   `json:"mode"`
	Error   bool     `json:"error"`
	Query   string   `json:"query"`
	Page    int      `json:"page"`
	Size    int      `json:"size"`
	Results []string `json:"results"`
}

func FofaSearchBody(keyword string) (hosts []string, err error) {
	qBase64 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("body=\"%s\"", keyword)))
	maxCount := GetConfig().Fofa.MaxCount
	if maxCount == 0 {
		maxCount = 100
	}
	url := fmt.Sprintf("https://fofa.info/api/v1/search/all?fields=host&full=true&qbase64=%s&email=%s&key=%s&size=%d",
		qBase64, GetConfig().Fofa.Email, GetConfig().Fofa.Key, maxCount)
	respBody, err := SimpleGet(url)
	if err != nil {
		return
	}
	var result FofaSearchResult
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return
	}

	if result.Error {
		err = errors.New(result.Errmsg)
		return
	}

	log.Printf("[*] FOFA found %d item, return %d hosts.\n", result.Size, len(result.Results))
	hosts = result.Results
	return
}

func FofaFetchContent(host string) (content string, err error) {
	url := fmt.Sprintf("https://fofa.info/result/website?host=%s", host)
	respBody, err := SimpleGet(url)
	if err != nil {
		return
	}

	return string(respBody), nil
}

func handleContent(content string) string {
	return html.UnescapeString(content)
}

func FofaMiner(format string) (links []string, err error) {
	config := GetConfig().Fofa
	log.Printf("[*] Starting discover from FOFA, threadCount=%d, maxCount=%d\n", config.ThreadCount, config.MaxCount)
	hosts, err := FofaSearchBody(format)
	if err != nil {
		return
	}
	log.Printf("[*] FOFA found %d host.\n", len(hosts))

	var mu sync.Mutex
	makeF := func(host string) func() {
		return func() {
			content, err := FofaFetchContent(host)
			if err != nil {
				log.Printf("[*] (%s) fetch content failed as %s.\n", host, err)
			} else {
				content = handleContent(content)
				l := MatchLinks(content, format)
				log.Printf("[*] (%s) fetch content size %d, matched %d links\n", host, len(content), len(l))
				mu.Lock()
				defer mu.Unlock()
				links = append(links, l...)
			}
		}
	}
	pool := simplethreadpool.NewSimpleThreadPool(config.ThreadCount)
	for _, host := range hosts {
		pool.Put(makeF(host))
	}
	pool.Sync()

	links = RemoveRep(links)
	log.Printf("[*] FOFA found %d link.\n", len(links))
	return
}
