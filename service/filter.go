package service

import (
	"github.com/spaolacci/murmur3"
	"hash"
	"math"
	"sync"
)

type BloomFilter struct {
	m       uint
	k       uint
	buckets []bool
	hashfn  []hash.Hash32
	mux     sync.Mutex
}

func NewBloomFilter(m uint, p float64) *BloomFilter {
	k := uint(math.Ceil(-math.Log2(p)))
	hashfn := make([]hash.Hash32, k)
	for i := uint(0); i < k; i++ {
		hashfn[i] = murmur3.New32WithSeed(uint32(i))
	}
	return &BloomFilter{
		m:       m,
		k:       k,
		buckets: make([]bool, m),
		hashfn:  hashfn,
	}
}

func (bf *BloomFilter) Add(key string) {
	bf.mux.Lock()
	defer bf.mux.Unlock()
	for i := uint(0); i < bf.k; i++ {
		hash := bf.hashfn[i]
		hash.Write([]byte(key))
		index := uint(hash.Sum32()) % bf.m
		bf.buckets[index] = true
		hash.Reset()
	}
}

func (bf *BloomFilter) Check(key string) bool {
	bf.mux.Lock()
	defer bf.mux.Unlock()
	for i := uint(0); i < bf.k; i++ {
		hash := bf.hashfn[i]
		hash.Write([]byte(key))
		index := uint(hash.Sum32()) % bf.m
		if !bf.buckets[index] {
			return false
		}
		hash.Reset()
	}
	return true
}
