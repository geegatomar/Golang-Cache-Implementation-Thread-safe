package evictorFirstN

// This eviction strategy is to show that we can easily plug-in any new eviction strategy easily,
// by only implementing the Evictor interface

// EvictorFirstN implements the Evictor interface
type EvictorFirstN struct {
	keys []string
}

func (e *EvictorFirstN) Evict(evictElement func(key string) bool) {
	for {
		ele := e.keys[len(e.keys)-1]
		e.keys = e.keys[:len(e.keys)-1]
		if !evictElement(ele) {
			return
		}
	}
}

func (e *EvictorFirstN) OnAdd(key string) error {
	if e.keys == nil {
		e.keys = make([]string, 0)
	}
	e.keys = append(e.keys, key)
	return nil
}

func (e *EvictorFirstN) OnAccess(key string) error {
	return nil
}

func (e *EvictorFirstN) OnDelete(key string) error {
	for ind, ele := range e.keys {
		if ele == key {
			e.keys = append(e.keys[:ind], e.keys[ind+1:]...)
			return nil
		}
	}
	return nil
}

func (e *EvictorFirstN) OnUpdate(key string, value string) error {
	return nil
}
