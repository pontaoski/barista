package discord

import (
	lru "github.com/hashicorp/golang-lru"
)

var commandCache *lru.ARCCache
var tagCache *lru.ARCCache

func init() {
	var err error
	commandCache, err = lru.NewARC(1 << 16)
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	tagCache, err = lru.NewARC(1 << 16)
	if err != nil {
		panic(err)
	}
}
