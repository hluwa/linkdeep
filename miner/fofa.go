package miner

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
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
	//qBase64 := base64.StdEncoding.EncodeToString([]byte("body=\"kakaotalk://shopping\" || body=\"kakaotalk://mart\" || body=\"kakaoopen://join\""))
	maxCount := GetConfig().Fofa.MaxCount
	if maxCount == 0 {
		maxCount = 100
	}
	url := fmt.Sprintf("https://fofa.so/api/v1/search/all?fields=host&full=true&qbase64=%s&email=%s&key=%s&size=%d",
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
	url := fmt.Sprintf("https://fofa.so/result/website?host=%s", host)
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

	channel := make(chan []string, config.ThreadCount)
	for _, host := range hosts {
		go func(h string) {
			content, err := FofaFetchContent(h)
			if err != nil {
				log.Printf("[*] (%s) fetch content failed as %s.\n", h, err)
				channel <- nil
			} else {
				content = handleContent(content)
				l := MatchLinks(content, format)
				log.Printf("[*] (%s) fetch content size %d, matched %d links\n", h, len(content), len(l))
				channel <- l
			}
		}(host)
	}

	for i := 0; i < len(hosts); i++ {
		l := <-channel
		links = append(links, l...)
	}
	links = RemoveRep(links)
	log.Printf("[*] FOFA found %d link.\n", len(links))
	return
}
