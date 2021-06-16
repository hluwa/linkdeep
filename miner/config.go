package miner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var config *Config
var mu sync.Mutex

type Config struct {
	Fofa struct {
		Email       string `json:"email"`
		Key         string `json:"key"`
		MaxCount    int    `json:"maxCount"`
		ThreadCount int    `json:"threadCount"`
		Proxy       string `json:"proxy"`
	} `json:"fofa"`

	Proxy string `json:"proxy"`
}

func (c *Config) GetFofaProxy() string {
	if c.Fofa.Proxy != "" {
		return c.Fofa.Proxy
	} else {
		return c.Proxy
	}
}

var ConfigPath = "config.json"

func GetConfig() *Config {
	mu.Lock()
	defer mu.Unlock()
	if config == nil {
		content, err := ioutil.ReadFile(ConfigPath)
		if err != nil {
			panic(fmt.Sprintf("Unable read config.json: %s", err))
		}
		var conf Config
		err = json.Unmarshal(content, &conf)
		if err != nil {
			panic(fmt.Sprintf("Unable unmarshal config: %s", err))
		}
		config = &conf

		if config.Fofa.MaxCount <= 0 {
			config.Fofa.MaxCount = 100
		} else if config.Fofa.MaxCount > 10000 {
			config.Fofa.MaxCount = 10000
		}
		if config.Fofa.ThreadCount <= 0 {
			config.Fofa.ThreadCount = 20
		}
	}
	return config
}
