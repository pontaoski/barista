package telegram

import lru "github.com/hashicorp/golang-lru"

var paginatorCache *lru.ARCCache

func init() {
	var err error
	paginatorCache, err = lru.NewARC(1 << 16)
	if err != nil {
		panic(err)
	}
}
