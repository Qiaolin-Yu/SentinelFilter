package service

import (
	"encoding/json"
	"fmt"
	"github.com/spaolacci/murmur3"
	"hash"
	"math"
	"net/http"
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

func (bf *BloomFilter) add(key string) {
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

func (bf *BloomFilter) check(key string) bool {
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

func (bf *BloomFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bf.handleGet(w, r)
	case http.MethodPost:
		bf.handlePost(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func (bf *BloomFilter) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}
	ok := bf.check(key)
	fmt.Fprint(w, ok)
}

func (bf *BloomFilter) handlePost(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Key string `json:"key"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if data.Key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}
	bf.add(data.Key)
	w.WriteHeader(http.StatusNoContent)
}
