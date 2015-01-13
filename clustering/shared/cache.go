package shared

import "sync"

// A Cacher just stores values for an (a,b).
// The values are symmetric. Values must be positive, returns -1 if not found.
type Cacher interface {
	Get(a, b int) float64
	Set(a, b int, value float64)
}

// A MapCache caches the values in a 2d hash map.
// bad performance for small sets, but scales better to large amounts of data.
type MapCache struct {
	sync.RWMutex
	count int
	cache map[int]map[int]float64
}

// NewMapCache creates a initializes a new map cache.
func NewMapCache(count int) *MapCache {
	return &MapCache{
		count: count,
		cache: make(map[int]map[int]float64, count),
	}
}

// Get returns the value, or -1 if not set.
func (c *MapCache) Get(a, b int) float64 {
	if a > b {
		a, b = b, a
	}

	c.RLock()
	defer c.RUnlock()

	if _, ok := c.cache[a]; !ok {
		return -1
	}

	if _, ok := c.cache[a][b]; !ok {
		return -1
	}

	return c.cache[a][b]
}

// Set sets a specific value.
func (c *MapCache) Set(a, b int, value float64) {
	if a > b {
		a, b = b, a
	}

	c.Lock()
	defer c.Unlock()

	if _, ok := c.cache[a]; !ok {
		c.cache[a] = make(map[int]float64, c.count/2)
	}

	c.cache[a][b] = value
}

// ArrayCache stores the values in a 2d array. This can cause issues
// if the data is very large.
type ArrayCache struct {
	count int
	cache []float64
}

// NewArrayCache initializes a new array cache
func NewArrayCache(count int) *ArrayCache {
	// TODO: figure out how to make this half the size since the values are symmetric.
	cache := make([]float64, count*count)
	for i := range cache {
		cache[i] = -1
	}

	return &ArrayCache{
		count: count,
		cache: cache,
	}
}

// Get returns the value, or -1 if not set.
func (c *ArrayCache) Get(a, b int) float64 {
	if a > b {
		a, b = b, a
	}

	return c.cache[a*c.count+b]
}

// Set sets the value.
func (c *ArrayCache) Set(a, b int, value float64) {
	if a > b {
		a, b = b, a
	}

	c.cache[a*c.count+b] = value
}
