package cache

import "time"

type Element struct {
	key            string
	value          string
	timeOfCreation time.Time
	// timeOfLastAccess time.Time   ---> Not using this anymore in cache since we now do this in the evictor itself
}

func (e *Element) Init(key string, val string) {
	e.key = key
	e.value = val
	e.timeOfCreation = time.Now()
}
