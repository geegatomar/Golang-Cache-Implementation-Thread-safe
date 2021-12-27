package evictor

// Evictor is an interface which any eviction strategy we want to use must implement this
type Evictor interface {
	Evict(func(key string) bool)
	OnAdd(key string) error
	OnAccess(key string) error
	OnDelete(key string) error
	OnUpdate(key string, value string) error
}
