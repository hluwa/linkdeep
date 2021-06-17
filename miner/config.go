package miner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var config *Config
var mu sync.Mutex

type MinerConfig struct {
	MaxCount    int    `json:"maxCount"`
	ThreadCount int    `json:"threadCount"`
	Proxy       string `json:"proxy"`
}

type Config struct {
	Fofa struct {
		MinerConfig
		Email string `json:"email"`
		Key   string `json:"key"`
	} `json:"fofa"`
	Github struct {
		MinerConfig
		Token string `json:"token"`
	} `json:"github"`
	Proxy string `json:"proxy"`
}

func (c *Config) GetFofaProxy() string {
	if c.Fofa.Proxy != "" {
		return c.Fofa.Proxy
	} else {
		return c.Proxy
	}
}

func (c *Config) GetGithubProxy() string {
	if c.Github.Proxy != "" {
		return c.Github.Proxy
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
		config = NewConfig()
		err = json.Unmarshal(content, config)
		if err != nil {
			panic(fmt.Sprintf("Unable unmarshal config: %s", err))
		}

		if config.Fofa.MaxCount > 10000 {
			config.Fofa.MaxCount = 10000
		}
	}
	return config
}

func NewConfig() *Config {
	return &Config{
		Fofa: struct {
			MinerConfig
			Email string `json:"email"`
			Key   string `json:"key"`
		}{
			MinerConfig: MinerConfig{
				MaxCount:    100,
				ThreadCount: 5,
			},
		},
		Github: struct {
			MinerConfig
			Token string `json:"token"`
		}{
			MinerConfig: MinerConfig{
				MaxCount:    1000,
				ThreadCount: 20,
			},
		},
	}
}
