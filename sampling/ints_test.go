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

package sampling_test

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/greatroar/randstat/sampling"

	"github.com/stretchr/testify/assert"
)

func testInts(t *testing.T, ints func(int, int64, rand.Source64) []int64) {
	t.Helper()

	for _, samplesize := range []int{1, 2, 3, 19, 10000} {
		for _, population := range []int64{1, 2, 5, 123, 1007, 9999999} {
			name := fmt.Sprintf("%d %d", samplesize, population)
			t.Run(name, func(t *testing.T) {
				sample := ints(samplesize, population, nil)
				checkSample(t, sample, samplesize, population)
			})
		}
	}
}

func toInt64(slice interface{}) []int64 {
	v := reflect.ValueOf(slice)
	out := make([]int64, 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		out = append(out, v.Index(i).Int())
	}
	return out
}

func TestInts32(t *testing.T) {
	t.Parallel()

	testInts(t, func(k int, n int64, r rand.Source64) []int64 {
		s := sampling.Ints32(k, int32(n), r, nil)
		return toInt64(s)
	})
}

func TestInts(t *testing.T) {
	t.Parallel()

	testInts(t, func(k int, n int64, r rand.Source64) []int64 {
		s := sampling.Ints(k, int(n), r, nil)
		return toInt64(s)
	})
}

func TestInts64(t *testing.T) {
	t.Parallel()

	testInts(t, func(k int, n int64, r rand.Source64) []int64 {
		return sampling.Ints64(k, n, r, nil)
	})
}

// Quick statistical test.
func TestInts64Stats(t *testing.T) {
	t.Parallel()

	const (
		population = 100
		samplesize = population / 10
	)

	freq := make([]float64, population)
	sample := make([]int64, 0, samplesize)
	r := rand.NewSource(0x52109).(rand.Source64)

	for i := 0; i < population*samplesize; i++ {
		sample = sampling.Ints64(samplesize, population, r, sample[:0])
		for _, x := range sample {
			freq[x]++
		}
	}

	var errNorm float64
	for _, f := range freq {
		err := math.Abs(f-population) / population
		errNorm += err * err
	}
	errNorm = math.Sqrt(errNorm) / float64(len(freq))
	assert.Less(t, errNorm, .015)
}

func TestIntsZero(t *testing.T) {
	buf := make([]int, 4)
	sample := sampling.Ints(0, 0xffff, nil, buf)
	assert.Equal(t, buf, sample)
	assert.Same(t, &buf[0], &sample[0])

	buf32 := make([]int32, 4)
	sample32 := sampling.Ints32(0, 0xfffffff, nil, buf32)
	assert.Equal(t, buf32, sample32)
	assert.Same(t, &buf32[0], &sample32[0])

	buf64 := make([]int64, 4)
	sample64 := sampling.Ints64(0, 0xfffffff, nil, buf64)
	assert.Equal(t, buf64, sample64)
	assert.Same(t, &buf64[0], &sample64[0])
}

func checkSample(t *testing.T, sample []int64, samplesize int, population int64) {
	t.Helper()
	assert := assert.New(t)

	if int64(samplesize) > population {
		samplesize = int(population)
	}
	assert.Equal(samplesize, len(sample))

	set := make(map[int64]struct{})
	for _, x := range sample {
		assert.GreaterOrEqual(x, int64(0))
		assert.Less(x, population)
		set[x] = struct{}{}
	}
	assert.Equal(samplesize, len(set))
}

func benchmarkInts64(b *testing.B, s int, n int64) {
	sample := make([]int64, s)
	r := rand.NewSource(time.Now().UnixNano()).(rand.Source64)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sample = sampling.Ints64(s, n, r, sample)
	}
}

func BenchmarkInts64_10_1e5(b *testing.B)   { benchmarkInts64(b, 10, 1e5) }
func BenchmarkInts64_10_1e6(b *testing.B)   { benchmarkInts64(b, 10, 1e6) }
func BenchmarkInts64_100_1e6(b *testing.B)  { benchmarkInts64(b, 100, 1e6) }
func BenchmarkInts64_1000_1e6(b *testing.B) { benchmarkInts64(b, 1000, 1e6) }
func BenchmarkInts64_1000_1e7(b *testing.B) { benchmarkInts64(b, 1000, 1e7) }
