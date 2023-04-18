package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
	bf := NewBloomFilter(10000, 0.01)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bf.Add(fmt.Sprintf("key%d", i))
			if i%2 != 0 {
				assert.True(t, bf.Check(fmt.Sprintf("key%d", i)))
			}
		}(i)
	}
	wg.Wait()
}

func TestFilter(t *testing.T) {
	bf := NewBloomFilter(1000, 0.01)
	bf.Add("hello")
	bf.Add("world")
	fmt.Println(bf.Check("hello")) // true
	fmt.Println(bf.Check("world")) // true
	fmt.Println(bf.Check("foo"))   // false
	fmt.Println(bf.Check("a"))     // false
}
