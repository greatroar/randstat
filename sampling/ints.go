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

package sampling

import (
	"math"
	"math/rand"

	"github.com/greatroar/randstat"
	"github.com/greatroar/randstat/internal/source"
)

// Ints appends to buf a simple random sample of the integers [0,n)
// and returns the resulting slice. The sample is not sorted.
//
// Random numbers are taken from r, or math.rand's global RNG if r is nil.
//
// If samplesize > n, the sample will be of size n instead.
//
// The time complexity of this function is O(s(1+log(n/s))).
func Ints(samplesize int, n int64, r rand.Source64, buf []int64) []int64 {
	if r == nil {
		r = source.Std
	}
	if int64(samplesize) > n {
		samplesize = int(n)
	}
	if samplesize == 0 {
		return buf
	}

	// Algorithm L from Li 1994, Reservoir-Sampling Algorithms of Time
	// Complexity O(n(1+log(N/n))), ACM TOMS,
	// https://doi.org/10.1145%2F198429.198435
	for i := 0; i < samplesize; i++ {
		buf = append(buf, int64(i))
	}
	sample := buf[len(buf)-samplesize:]

	var (
		w = float64(1)
		i = float64(samplesize)
		k = float64(samplesize)
		N = float64(n)
	)
	for {
		w *= math.Exp(math.Log(random01(r)) / k)
		i += 1 + math.Floor(math.Log(random01(r))/math.Log1p(-w))
		if i >= N {
			break
		}
		j := randstat.Intn(r, len(sample))
		sample[j] = int64(i)
	}

	return buf
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
