package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/geegatomar/cacheenterpret/cache"
	"github.com/stretchr/testify/assert"
)

func printCacheElements(c *cache.Cache) {
	fmt.Println("====================================================================")
	c.ViewCacheElements()
	fmt.Println("====================================================================")
}

func otherOperationsOnCache(c *cache.Cache, wg *sync.WaitGroup) {
	defer wg.Done()
	c.Put("9", "9999")
	c.Put("10", "10000")
	c.Put("11", "110000")
	fmt.Println(c.Get("9"))
	c.Delete("9")
	fmt.Println(c.Get("9"))

	for i := 100; i < 3100; i++ {
		c.Put(fmt.Sprint(i), fmt.Sprint(i))
		if i%4 == 0 {
			c.Get(fmt.Sprint(i))
			c.Delete(fmt.Sprint(i))
			c.Get(fmt.Sprint(i))
		}
	}
	fmt.Println("Going to delete 8...")
	c.Delete("8")

}

// This test is for Stress testing where we have mutliple goroutines running simultaneously and
// putting/getting/deleting multiple values from the cache which is used to test that the code
// is thread-safe and does not go into any deadlocks.
func TestStress(t *testing.T) {
	// Using wait groups since we want to make use of multiple goroutines that we need to wait for
	// before main finishes execution.
	wg := new(sync.WaitGroup)

	c := new(cache.Cache)
	c.Init(5, cache.LRU)
	c.Put("1", "1111")
	c.Put("2", "2222")
	c.Put("3", "3333")
	c.Put("4", "4444")
	c.Put("5", "5555")
	fmt.Println(c.Get("1"))
	c.Put("6", "6666")
	fmt.Println(c.Get("4"))
	c.Put("7", "7777")
	c.Put("8", "8888")
	fmt.Println(c.GetSize())

	// Calling other goroutine here
	fmt.Println("NEW GOROUTINE HERE")
	wg.Add(2)
	go otherOperationsOnCache(c, wg)
	go otherOperationsOnCache(c, wg)

	fmt.Println(c.Get("6"))
	fmt.Println(c.Get("5"))
	fmt.Println(c.Get("7"))
	c.Delete("7")
	fmt.Println(c.Get("7"))
	c.Put("8", "9999999999")

	for i := 500; i < 1520; i++ {
		c.Put(fmt.Sprint(i), fmt.Sprint(i))
		if i%4 == 0 {
			c.Get(fmt.Sprint(i))
			c.Delete(fmt.Sprint(i))
			c.Get(fmt.Sprint(i))
		}
	}

	wg.Wait()

	printCacheElements(c)
	fmt.Println(c.GetSize())

	assert.Equal(t, true, c.GetSize() <= 5, "Asserting size of cache is less than the max value")
}
