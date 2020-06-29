package sampling_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/greatroar/randstat"
	"github.com/greatroar/randstat/sampling"
	"github.com/greatroar/randstat/xoshiro256"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaroptBasic(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	const size = 10

	sample := sampling.NewVaropt(size, rand.NewSource(42).(rand.Source64))
	assert.Equal(0, sample.Len())
	assert.Panics(func() { sample.Item(0) })

	var all, rejected []interface{}
	show := func(x interface{}, w float64) interface{} {
		require.Greater(t, w, 0.) // Sanity check.

		reject := sample.Show(x, w)
		all = append(all, x)
		if reject != nil {
			rejected = append(rejected, reject)
		}
		return reject
	}

	show0 := func(x interface{}) {
		n := sample.Len()
		reject := sample.Show(x, 0) // Zero weight, don't accept.
		require.Equal(t, x, reject)
		require.Equal(t, n, sample.Len())
	}

	show0("no")

	reject := show("yes", 1)
	assert.Equal(1, sample.Len())
	assert.Nil(reject)

	for i := 0; i < size; i++ {
		show0("still not")
		reject = show(i, float64(1+i))
	}
	assert.Equal(size, sample.Len())
	assert.NotNil(reject)

	for i := 0; i < 2*size; i++ {
		show0("nope")
		reject = show(i, float64(1+i))
		assert.Equal(size, sample.Len())
		assert.NotNil(reject)
	}

	sampled := make([]interface{}, sample.Len())
	for i := range sampled {
		sampled[i] = sample.Item(i)
	}
	assert.ElementsMatch(all, append(rejected, sampled...))
}

// Quick statistical test.
func TestVaroptStats(t *testing.T) {
	t.Parallel()

	const (
		population = 2e6
		samplesize = 2e4
	)

	prob := [...]float64{.01, .1, .2, .5}
	s := sampling.NewVaropt(samplesize, rand.NewSource(42).(rand.Source64))

	for i := 0; i < population; i++ {
		s.Show(i, prob[i%len(prob)])
	}

	assert.EqualValues(t, samplesize, s.Len())

	freq := [len(prob)]float64{}
	for i := 0; i < s.Len(); i++ {
		freq[s.Item(i).(int)%len(prob)]++
	}

	var errNorm float64
	for i, f := range freq {
		exp := samplesize * prob[i]

		err := math.Abs(f-exp) / exp
		errNorm += err * err
	}

	errNorm = math.Sqrt(errNorm) / float64(len(freq))
	assert.Less(t, errNorm, .12)
}

func benchmarkVaropt(b *testing.B, k, n int) {
	r := xoshiro256.New(uint64(time.Now().UnixNano()))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sample := sampling.NewVaropt(k, r)
		for j := 0; j < n; j++ {
			sample.Show(nil, randstat.Float64(r))
		}
	}
}

func BenchmarkVaropt10_1e5(b *testing.B)   { benchmarkVaropt(b, 10, 1e5) }
func BenchmarkVaropt10_1e6(b *testing.B)   { benchmarkVaropt(b, 10, 1e6) }
func BenchmarkVaropt100_1e6(b *testing.B)  { benchmarkVaropt(b, 100, 1e6) }
func BenchmarkVaropt10_1e7(b *testing.B)   { benchmarkVaropt(b, 10, 1e7) }
func BenchmarkVaropt100_1e7(b *testing.B)  { benchmarkVaropt(b, 100, 1e7) }
func BenchmarkVaropt1000_1e7(b *testing.B) { benchmarkVaropt(b, 1000, 1e7) }
