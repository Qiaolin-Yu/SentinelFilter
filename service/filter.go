package service

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"net/http"
	"sync"
)

type BloomFilter struct {
	bitArray []bool
	k        uint32 // number of hash functions
	m        uint32 // number of bits in the filter
	mu       sync.Mutex
}

func NewBloomFilter(n uint32, p float64) *BloomFilter {
	m := uint32(-1 * float64(n) * math.Log(p) / (math.Ln2 * math.Ln2))
	k := uint32(math.Ceil((float64(m) / float64(n)) * math.Ln2))
	return &BloomFilter{make([]bool, m), k, m, sync.Mutex{}}
}

func (bf *BloomFilter) add(value string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	for i := uint32(0); i < bf.k; i++ {
		hash := fnv.New32a()
		hash.Write([]byte(fmt.Sprintf("%s%d", value, i)))
		index := hash.Sum32() % bf.m
		bf.bitArray[index] = true
	}
}

func (bf *BloomFilter) contains(value string) bool {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	for i := uint32(0); i < bf.k; i++ {
		hash := fnv.New32a()
		hash.Write([]byte(fmt.Sprintf("%s%d", value, i)))
		index := hash.Sum32() % bf.m
		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) delete(value string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()
	for i := uint32(0); i < bf.k; i++ {
		hash := fnv.New32a()
		hash.Write([]byte(fmt.Sprintf("%s%d", value, i)))
		index := hash.Sum32() % bf.m
		bf.bitArray[index] = false
	}
}

func (bf *BloomFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bf.handleGet(w, r)
	case http.MethodPost:
		bf.handlePost(w, r)
	case http.MethodDelete:
		bf.handleDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (bf *BloomFilter) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}
	if bf.contains(key) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "true")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "false")
	}
}

func (bf *BloomFilter) handlePost(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Key string `json:"key"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if data.Key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}
	bf.add(data.Key)
	w.WriteHeader(http.StatusCreated)
}

func (bf *BloomFilter) handleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}
	bf.delete(key)
	w.WriteHeader(http.StatusNoContent)
}
