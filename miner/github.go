package miner

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hluwa/simplethreadpool"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type GithubAuthResponse struct {
	Message string `json:"message"`
}

type GithubSearchResponse struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []struct {
		Url string `json:"url"`
	} `json:"items"`
}

type GithubContentResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func GithubGet(u string) ([]byte, error) {
	return CustomGet(u, func(r *http.Request) {
		r.Header.Set("Authorization", fmt.Sprintf("token %s", GetConfig().Github.Token))
	})
}

func GithubAuth() error {
	content, err := GithubGet("https://api.github.com")
	if err != nil {
		return err
	}
	var resp GithubAuthResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
		return err
	}

	if resp.Message != "" {
		log.Printf("[*] Github auth failed at %s\n", resp.Message)
		return err
	}

	return nil
}

func GithubSearch(query string) (result []string, err error) {
	page := 0
	failed := 0
	for ; len(result) < GetConfig().Github.MaxCount; page++ {
		perPage := GetConfig().Github.MaxCount - len(result)
		if perPage > 100 {
			perPage = 100
		}
		params := url.Values{}
		params.Set("q", query)
		params.Set("per_page", strconv.Itoa(perPage))
		params.Set("page", strconv.Itoa(page))
		params.Set("sort", "")
		content, err := GithubGet(fmt.Sprintf("https://api.github.com/search/code?%s", params.Encode()))
		if err != nil {
			if failed > 5 {
				break
			} else {
				failed++
				continue
			}

		}
		var resp GithubSearchResponse
		err = json.Unmarshal(content, &resp)
		if err != nil || len(resp.Items) <= 0 {
			break
		}

		for _, item := range resp.Items {
			result = append(result, item.Url)
		}

		if resp.IncompleteResults {
			break
		}
	}
	return result, err
}

func GithubGetContent(url string) (content string, err error) {
	body, err := GithubGet(url)
	if err != nil {
		return
	}

	var resp GithubContentResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return
	}

	if resp.Content == "" {
		return "", errors.New(fmt.Sprintf("cannot fetch content, body at %s", string(body)))
	}

	if resp.Encoding != "base64" {
		return "", errors.New(fmt.Sprintf("unimplemented encoding at %s", resp.Encoding))
	}

	body, err = base64.StdEncoding.DecodeString(resp.Content)
	if err != nil {
		return
	}

	return string(body), nil
}

func GithubMiner(format string) (links []string, err error) {
	if err = GithubAuth(); err != nil {
		return nil, errors.New(fmt.Sprintf("github is cannot auth at %s", err))
	}
	config := GetConfig().Github
	log.Printf("[*] Starting discover from Github, threadCount=%d, maxCount=%d\n", config.ThreadCount, config.MaxCount)
	urls, err := GithubSearch(fmt.Sprintf("\"%s\"", format))
	if err != nil {
		return
	}

	var mu sync.Mutex
	makeF := func(u string) func() {
		return func() {
			content, err := GithubGetContent(u)
			if err != nil {
				log.Printf("[*] (%s) fetch content failed as %s.\n", u, err)
			} else {
				l := MatchLinks(content, format)
				log.Printf("[*] (%s) fetch content size %d, matched %d links\n", u, len(content), len(l))
				mu.Lock()
				defer mu.Unlock()
				links = append(links, l...)
			}
		}
	}
	pool := simplethreadpool.NewSimpleThreadPool(GetConfig().Github.ThreadCount)
	for _, u := range urls {
		pool.Put(makeF(u))
	}
	pool.Sync()
	links = RemoveRep(links)
	log.Printf("[*] Github found %d link.\n", len(links))
	return
}
