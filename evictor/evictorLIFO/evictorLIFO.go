package evictorLIFO

// EvictorLIFO implements the Evictor interface
type EvictorLIFO struct {
	keys []string
}

func (e *EvictorLIFO) Evict(evictElement func(key string) bool) {
	for {
		ele := e.keys[len(e.keys)-1]
		e.keys = e.keys[:len(e.keys)-1]
		if !evictElement(ele) {
			return
		}
	}
}

func (e *EvictorLIFO) OnAdd(key string) error {
	if e.keys == nil {
		e.keys = make([]string, 0)
	}
	e.keys = append(e.keys, key)
	return nil
}

func (e *EvictorLIFO) OnAccess(key string) error {
	return nil
}

func (e *EvictorLIFO) OnDelete(key string) error {
	for ind, ele := range e.keys {
		if ele == key {
			e.keys = append(e.keys[:ind], e.keys[ind+1:]...)
			return nil
		}
	}
	return nil
}

func (e *EvictorLIFO) OnUpdate(key string, value string) error {
	return nil
}
