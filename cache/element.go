package cache

import "time"

type Element struct {
	key            string
	value          string
	timeOfCreation time.Time
	// timeOfLastAccess time.Time
}

func (e *Element) Init(key string, val string) {
	e.key = key
	e.value = val
	e.timeOfCreation = time.Now()
}
