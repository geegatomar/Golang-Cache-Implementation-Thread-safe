package main

import (
	"testing"

	"github.com/geegatomar/cacheenterpret/cache"
	"github.com/stretchr/testify/assert"
)

// Some of the basic unittests to show the functionality of the cache created

func TestPutAndGetCache(t *testing.T) {
	c := new(cache.Cache)
	c.Init(5, cache.LRU)
	c.Put("1", "1111")
	c.Put("2", "2222")

	actual, _ := c.Get("1")
	expected := "1111"
	assert.Equal(t, expected, actual, "Asserting the returned element by cache")
}

func TestUpdatedPutCache(t *testing.T) {
	c := new(cache.Cache)
	c.Init(5, cache.LRU)
	c.Put("1", "1111")
	c.Put("2", "2222")
	c.Put("1", "3333")
	c.Put("1", "4444")
	// If we put/insert the same element multiple times, then we expect the value for that key
	// to get updated each time, and eventually get the latest updated value.

	actual, _ := c.Get("1")
	expected := "4444"
	assert.Equal(t, expected, actual, "Asserting the returned element by cache")
}

func TestFIFOEviction(t *testing.T) {
	c := new(cache.Cache)
	c.Init(3, cache.FIFO)
	c.Put("1", "1111")
	c.Put("2", "2222")
	c.Put("3", "3333")
	c.Get("1")
	c.Put("4", "4444")
	// Element with key "1" should be evicted after putting the 4th element, since max size
	// of the cache is 3

	actual := int(c.GetSize())
	expected := 3
	assert.Equal(t, expected, actual, "Asserting size of cache")

	actual_cache_elements := c.GetAllKeysSorted()
	expected_cache_elements := []string{"2", "3", "4"}
	assert.Equal(t, expected_cache_elements, actual_cache_elements, "Asserting elements of cache")
}

func TestLIFOEviction(t *testing.T) {
	c := new(cache.Cache)
	c.Init(3, cache.LIFO)
	c.Put("1", "1111")
	c.Put("2", "2222")
	c.Put("3", "3333")
	c.Get("1")
	c.Put("4", "4444")

	actual := int(c.GetSize())
	expected := 3
	assert.Equal(t, expected, actual, "Asserting size of cache")

	actual_cache_elements := c.GetAllKeysSorted()
	expected_cache_elements := []string{"1", "2", "3"}
	assert.Equal(t, expected_cache_elements, actual_cache_elements, "Asserting elements of cache")
}

func TestLRUEviction(t *testing.T) {
	c := new(cache.Cache)
	c.Init(3, cache.LRU)
	c.Put("1", "1111")
	c.Put("2", "2222")
	c.Put("3", "3333")
	c.Get("1")
	c.Put("2", "22222")
	c.Put("4", "4444")

	actual := int(c.GetSize())
	expected := 3
	assert.Equal(t, expected, actual, "Asserting size of cache")

	actual_cache_elements := c.GetAllKeysSorted()
	expected_cache_elements := []string{"1", "2", "4"}
	assert.Equal(t, expected_cache_elements, actual_cache_elements, "Asserting elements of cache")
}

func TestLRUDeletionFromCache(t *testing.T) {
	c := new(cache.Cache)
	c.Init(3, cache.LRU)
	c.Put("1", "1111")
	c.Put("2", "2222")
	c.Put("3", "3333")
	c.Get("1")

	// Before Deletion of element, the error returned is nil
	_, err := c.Get("3")
	assert.Equal(t, err, nil, "Asserting error returned by Get is nil")

	c.Delete("3")

	// After Deletion of element, the error returned is not nil since element was not found in cache
	_, err = c.Get("3")
	assert.NotEqual(t, err, nil, "Asserting error returned by Get is not nil")

	c.Put("4", "4444")
	c.Put("5", "5555")

	actual := int(c.GetSize())
	expected := 3
	assert.Equal(t, expected, actual, "Asserting size of cache")

	actual_cache_elements := c.GetAllKeysSorted()
	expected_cache_elements := []string{"1", "4", "5"}
	assert.Equal(t, expected_cache_elements, actual_cache_elements, "Asserting elements of cache")
}
