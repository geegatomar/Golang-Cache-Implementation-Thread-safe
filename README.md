# Cache Implementation in GoLang

### Problem Statement
Design and Implement an in-memory caching library for general use.
<br>
*Must Have*
- Support for multiple Standard Eviction Policies ( FIFO, LRU, LIFO ) <br>
- Support to add custom eviction policies
<br>
*Good To Have*
- Thread saftey

### My Solution

#### Explanation of the project
- TODO: Make is thread safe, and explain that here.

- TODO: Explain the LRU & heap(priority_queue)

- TODO: Talk abt the eviction percentage and why we set it (to avoid the costly operations on evict)



#### Highlights of the project
- The evictor has been made an interface, and every/any strategy (LRU, FIFO, LIFO, etc) we want to add-on can be done by implementing this interface.
- The cache implementation and the eviction strategies are decoupled, making it easily pluggable for anyone to add more strategies.
- Its is thread safe since we are making use of RWMutex and acquiring lock whenever we are accessing the shared memory in the implemented cache (in our case it is both c.kv and c.ev).
- The evictor interface need not take care of thread safety as the cache itself will take care of that.
- We have set an eviction percentage, hence Evict is not called very often, and whenever called, it is called for very few entries.


#### Simple Testing

#### Caveats and Future Improvements


