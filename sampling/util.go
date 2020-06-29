package sampling

import (
	"math/rand"

	"github.com/greatroar/randstat"
	"github.com/greatroar/randstat/xoshiro256"
)

func maybeXoshiro(r rand.Source64) rand.Source64 {
	if r == nil {
		r = xoshiro256.New(rand.Uint64())
	}
	return r
}

// random01 returns a random float64 in (0,1).
func random01(r rand.Source) float64 {
	// Go 1.14 will not inline the obvious loop, but it will inline this.
retry:
	x := randstat.Float64(r)
	if x == 0 {
		goto retry
	}
	return x
}
