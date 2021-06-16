package main

import (
	"./miner"
	"testing"
)

func TestFofa(t *testing.T) {
	links, e := miner.FofaMiner("tencent://")
	if e != nil {
		t.Fatal(e)
	}
	for _, links := range links {
		t.Log(links)
	}
}
