package evictorFIFO

// EvictorFIFO implements the Evictor interface
type EvictorFIFO struct {
	keys []string
}

func (e *EvictorFIFO) Evict(evictElement func(key string) bool) {
	for {
		ele := e.keys[0]
		e.keys = e.keys[1:len(e.keys)]
		if !evictElement(ele) {
			return
		}
	}
}

func (e *EvictorFIFO) OnAdd(key string) error {
	if e.keys == nil {
		e.keys = make([]string, 0)
	}
	e.keys = append(e.keys, key)
	return nil
}

func (e *EvictorFIFO) OnAccess(key string) error {
	return nil
}

func (e *EvictorFIFO) OnDelete(key string) error {
	for ind, ele := range e.keys {
		if ele == key {
			e.keys = append(e.keys[:ind], e.keys[ind+1:]...)
			return nil
		}
	}
	return nil
}

func (e *EvictorFIFO) OnUpdate(key string, value string) error {
	return nil
}
