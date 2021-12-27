# Cache Implementation in GoLang

## Problem Statement
Design and Implement an in-memory caching library for general use.
<br>
<br>
*Must Have*
- Support for multiple Standard Eviction Policies ( FIFO, LRU, LIFO ) <br>
- Support to add custom eviction policies
<br>

*Good To Have*
- Thread saftey

<br>

## My Solution

### Explanation of the project
- The *cache* package has the implementation of the main cache, which has a map of key-value pair as its member.
- There is another package called *evictor* which only handles the eviction, and is independent of the cache package; hence decoupling the logic of the cache and the eviction strategies.
- The eviction function is called inside the Put() method after inserting an element into the cache.
- We are setting an eviction percentage which decides how many of the current elements in the cache will get evicted. So if the current size of the cache is 100 (and has exceeded the max limit of the cache size), and the eviction percentage is 30%, then after the eviction only 70 entries will remain in the cache. Reason for doing so is that we don't want to call the expensive evict operation very frequently.
- The implementation for FIFO, LIFO and LRU have been included in their respective packages.
- For LRU's implementation, we are making use of the timeOfLastAccess which gets updated whenever the element was accessed. And we use a minHeap to implement it in our code using the PriorityQueue which always keeps the element with the oldest time on top for extraction next.


### Highlights of the project
- The evictor has been made an interface, and every/any strategy (LRU, FIFO, LIFO, Random, etc) we want to add-on can be done by implementing this interface.
- The cache implementation and the eviction strategies are decoupled, making it easily pluggable for anyone to add more strategies.
- Its is thread safe since we are making use of RWMutex and acquiring lock whenever we are accessing the shared memory in the implemented cache (in our case it is both c.kv and c.ev).
- The evictor interface need not take care of thread safety as the cache itself will take care of that.
- We have set an eviction percentage, hence Evict is not called very often, and whenever called, it is called for few entries.


### Simple Testing
The two major test files are: 
1. main_stress_test.go 
  - This is doing an end to end stress test with multiple go routines where we are simultaneuously adding/getting/deleting from the cache. Hence also checks if its thread safe and ensures that we don't enter a deadlock.
  - To run:  `go test -v main_stress_test.go`
2. main_test.go.
  - This has muliple smaller individual unittests to test smaller functionalities in the code.
  - To run:  `go test -v main_test.go`
 


### Caveats and Future Improvements
1. Other data structures can be explored for implementing the LRU eviction. Currently we have used a PriorityQueue (minHeap), but other options can be explored here which give a better time complexity on some of the operations.
And hence the internal implementations of these functions and methods can be changed to modify the 
2. More fine grained mutexes can be used. 
3. Adding more extensive testing.
4. Ability to run the eviction as a background thread (similar to a daemon thread like a garbage collector) which always evicts once it detects that cache size is exceeding the max limit.
