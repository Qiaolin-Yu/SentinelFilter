package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bf.add(fmt.Sprintf("key%d", i))
			if i%2 == 0 {
				bf.delete(fmt.Sprintf("key%d", i))
			}
			if i%2 != 0 {
				assert.True(t, bf.contains(fmt.Sprintf("key%d", i)))
			} else {
				assert.False(t, bf.contains(fmt.Sprintf("key%d", i)))
			}
		}(i)
	}
	wg.Wait()
}

func TestFilter(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01)
	bf.add("hello")
	bf.add("world")
	fmt.Println(bf.contains("hello")) // true
	fmt.Println(bf.contains("world")) // true
	fmt.Println(bf.contains("foo"))   // false
	bf.delete("hello")
	fmt.Println(bf.contains("hello")) // false
}
