package simulation

import (
	"math/rand"
	"sync"
)

// RandomSource is the only randomness dependency used by the simulator.
type RandomSource interface {
	Float64() float64
	IntBetween(min int, max int) int
}

// SeededRandom is a deterministic RandomSource backed by math/rand.Rand.
type SeededRandom struct {
	mu sync.Mutex
	r  *rand.Rand
}

// NewSeededRandom creates a deterministic random source.
func NewSeededRandom(seed int64) *SeededRandom {
	return &SeededRandom{r: rand.New(rand.NewSource(seed))}
}

// Float64 returns a random float in [0, 1).
func (r *SeededRandom) Float64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.r.Float64()
}

// IntBetween returns an integer in the inclusive range [min, max].
func (r *SeededRandom) IntBetween(min int, max int) int {
	if max <= min {
		return min
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	return min + r.r.Intn(max-min+1)
}
