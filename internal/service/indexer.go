package service

import "sync/atomic"

type Indexer struct {
	val atomic.Int64
}

func NewIndexer() *Indexer {
	return &Indexer{
		val: atomic.Int64{},
	}
}

func (i *Indexer) GetVal() int64 {
	return i.val.Add(1)
}
