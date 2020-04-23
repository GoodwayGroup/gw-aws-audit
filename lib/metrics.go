package lib

import "sync"

type Metrics struct {
	Volumes int
	SumVolumeSize int64
	VolumeCosts int
	Snapshots int
	SumSnapshotSize int64
	Processed int64
	Skipped int64
	Modified int64
}

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	defer c.mux.Unlock()
	return c.v[key]
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key]++
	c.mux.Unlock()
}

// Add value the counter for the given key.
func (c *SafeCounter) Add(key string, value int) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key] += value
	c.mux.Unlock()
}