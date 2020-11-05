// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package randstat

import (
	"math/bits"
	"math/rand"
)

const (
	maxint32  = 1<<31 - 1
	maxuint32 = 1<<32 - 1
)

// Intn returns a uniformly random integer from the range [0,n).
//
// n must be greater than 0. r must not be nil.
//
// This function is designed to be faster than math/rand.Rand.Intn on 64-bit
// machines. It is likely slower on 32-bit machines.
func Intn(r rand.Source64, n int) int {
	// We unconditionally call Int63n, which gives a 45% speedup on amd64
	// compared to Int31n. We take the 50% performance hit in 386 mode
	// compared to the stdlib version (both when reading from a
	// math/rand.NewSource).
	return int(Int63n(r, int64(n)))
}

// Int31n returns a uniformly random integer from the range [0,n).
//
// n must be greater than 0. r must not be nil.
func Int31n(r rand.Source, n int32) int32 {
	// Algorithm (5) from https://arxiv.org/pdf/1805.10941.pdf.
	// We could check the popcnt of n for powers of two, but that slows us
	// down 5% for other numbers (on amd64).
	u := uint32(n)

	// 32-bit random number, as a uint64.
	r32 := func() uint64 {
		return uint64(r.Int63()) & maxuint32
	}

	m := r32() * uint64(u)
	lo := uint32(m)
	if lo < u {
		thresh := (-u) % u
		for lo < thresh {
			m = r32() * uint64(u)
			lo = uint32(m)
		}
	}

	return int32(m >> 32)
}

// Int63n returns a uniformly random integer from the range [0,n).
//
// n must be greater than 0. r must not be nil.
func Int63n(r rand.Source64, n int64) int64 {
	// See comment in Int31n.
	u := uint64(n)

	hi, lo := bits.Mul64(uint64(r.Uint64()), uint64(n))
	if lo < u {
		thresh := (-u) % u
		for lo < thresh {
			hi, lo = bits.Mul64(uint64(r.Uint64()), u)
		}
	}

	return int64(hi)
}

// Shuffle generates a random permutation.
//
// n must not be negative. r must not be nil.
func Shuffle(r rand.Source64, n int, swap func(i, j int)) {
	switch {
	case n == 0:
		return
	case n < 0:
		panic("randstat.Shuffle: n < 0")
	case r == nil:
		panic("randstat.Shuffle: no random source given")
	}

	// Fisher-Yates shuffle.
	for i := n - 1; i > maxint32-1; i-- {
		j := int(Int63n(r, int64(1+i)))
		if i != j {
			swap(i, j)
		}
	}

	for i := n - 1; i > 0; i-- {
		j := int(Int31n(r, int32(1+i)))
		if i != j {
			swap(i, j)
		}
	}
}
