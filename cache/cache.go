package cache

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/geegatomar/cacheenterpret/evictor"
	"github.com/geegatomar/cacheenterpret/evictor/evictorFIFO"
	"github.com/geegatomar/cacheenterpret/evictor/evictorFirstN"
	"github.com/geegatomar/cacheenterpret/evictor/evictorLIFO"
	"github.com/geegatomar/cacheenterpret/evictor/evictorLRU"
)

var (
	ErrorNotFound = errors.New("Error: Element not found in cache.")
)

const (
	// Setting the eviction percentage which decides how many of the total elements in
	// the cache it will evict once the cache limit exceeds
	EVICTION_PERCENTAGE = 30
)

// Eviction Strategies our cache supports (pluggable)
const (
	FIFO   = "FIFO"
	LRU    = "LRU"
	LIFO   = "LIFO"
	FirstN = "FirstN"
)

type Cache struct {
	// Key-value store where key is of type string, and the value is of type pointer to Element
	kv               map[string]*Element
	evictionStrategy string
	// ev is of type Evictor which is implemented by various eviction strategies and assigned accordingly for the cache
	ev           evictor.Evictor
	MaxSize      int32
	curSize      int32
	EvictPercent int32
	mutex        *sync.RWMutex
}

func (c *Cache) Init(size int32, evictionStrategy string) {
	c.MaxSize = size
	c.curSize = 0
	c.kv = make(map[string]*Element)
	c.evictionStrategy = evictionStrategy
	c.EvictPercent = EVICTION_PERCENTAGE
	c.mutex = new(sync.RWMutex)

	switch strategy := c.evictionStrategy; strategy {
	case FIFO:
		c.ev = &evictorFIFO.EvictorFIFO{}
	case LIFO:
		c.ev = &evictorLIFO.EvictorLIFO{}
	case LRU:
		c.ev = &evictorLRU.EvictorLRU{}
	case FirstN:
		c.ev = &evictorFirstN.EvictorFirstN{}
	default:
		panic("Invalid eviction strategy chosen")
	}
}

// Put inserts/sets a new key-value pair in the cache, and if the key already exists, then
// the value of that key gets updated to the latest value given.
func (c *Cache) Put(key string, value string) {
	e := new(Element)
	e.Init(key, value)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	// If the key is already present in the cache, then we only update it & dont increment current size
	if _, ok := c.kv[key]; ok {
		c.ev.OnUpdate(key, value)
	} else {
		c.curSize++
		c.ev.OnAdd(key)
	}

	c.kv[key] = e

	if c.curSize <= c.MaxSize {
		return
	}

	// Running an eviction whenever we add to the cache and the size of the cache exceeds the limit of the cache.
	func() {

		fmt.Println("Evicting now...  [Current size: ", c.curSize, " , Max size: ", c.MaxSize, "]")
		numEvicted := 0
		toEvict := int(c.EvictPercent * c.curSize / 100)
		fmt.Println("Evicting", toEvict, "out of", c.curSize, "elements...")
		fmt.Println("Size of cache: ", c.GetSize())

		// The Evict method takes a callback function which gets called and executed inside the
		// Evict method. Another way of doing this is to pass the entire cache as an argument in
		// the function but the down-side of that is that the cache implementation would become
		// tightly coupled but with this current implementation we have both the cache and evictor
		// decoupled and independent.

		// So here the internal Evict function will only return the key to be evicted each time and to
		// indicate when we should stop evicting, we have made this callback function return a boolean
		c.ev.Evict(func(key string) bool {
			numEvicted++
			// For each element selected by our evictor, we delete it from the map in the cache.
			fmt.Println("Inside cache.go, evicting element with key ", key)
			delete(c.kv, key)
			c.curSize--
			if numEvicted >= toEvict {
				return false
			}
			return true
		})
	}()
}

// Get returns the value of the specified key given if found, and returns the ErrorNotFound error
// with null string if key not found.
func (c *Cache) Get(key string) (string, error) {
	// The reason we took a Write lock here (and not a Read lock instead) is because we are modifying
	// in c.ev.OnAccess(key), and since we want to make the cache thread-safe and leave out the
	// evictor from having to handle all this, hence we handle everything here in the cache itself
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for k, v := range c.kv {
		if k == key {
			// c.kv[key].timeOfLastAccess = time.Now()   ---> All this logic has now been moved
			// into the evictor, so we dont need to handle it here.
			c.ev.OnAccess(key)
			return v.value, nil
		}
	}
	return "", ErrorNotFound
}

func (c *Cache) ViewCacheElements() {
	c.mutex.RLock()
	for k, v := range c.kv {
		fmt.Println("Element in cache: ", k, v)
	}
	c.mutex.RUnlock()
}

func (c *Cache) GetAllKeysSorted() []string {
	sorted := []string{}
	c.mutex.RLock()
	for k, v := range c.kv {
		sorted = append(sorted, k)
		fmt.Println("Element in cache: ", k, v)
	}
	sort.Strings(sorted)
	c.mutex.RUnlock()
	return sorted
}

func (c *Cache) GetSize() int32 {
	return c.curSize
}

func (c *Cache) Delete(key string) {
	// Calling Delete method on the evictor as well to update it too.
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.kv[key]; !ok {
		fmt.Println("Element to be deleted", key, "not found in cache.")
		return
	}
	c.curSize--
	c.ev.OnDelete(key)
	fmt.Println("Element", key, "successfully deleted.")
	delete(c.kv, key)
}
