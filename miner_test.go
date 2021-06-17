package main

import (
	"./miner"
	"fmt"
	"testing"
)

func Mining(name string, f func(string) ([]string, error), t *testing.T) ([]string, error) {
	links, e := f("tencent://message/")
	if e != nil {
		t.Fatal(e)
	}

	if len(links) <= 0 {
		t.Fatal(fmt.Sprintf("Unable discover from %s.", name))
	}
	return links, e
}

func TestFofa(t *testing.T) {
	_, _ = Mining("FOFA", miner.FofaMiner, t)
	//for _, links := range links {
	//	t.Log(links)
	//}
}

func TestGithub(t *testing.T) {
	_, _ = Mining("Github", miner.GithubMiner, t)
	//for _, links := range links {
	//	t.Log(links)
	//}
}
