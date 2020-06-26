// +build amd64 arm64 ppc64

package sampling_test

import (
	"testing"

	"github.com/greatroar/randstat/sampling"
)

// Test sampling from population > 1<<31 with Ints.
func TestInts64bit(t *testing.T) {
	const largepop = 1<<48 + 7
	sample := sampling.Ints(10, largepop, nil, nil)
	checkSample(t, toInt64(sample), 10, largepop)
}
