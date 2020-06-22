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

package randstat_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/greatroar/randstat"
)

type benchmarkSource struct {
	r      rand.Source64
	ncalls int64
}

func (s *benchmarkSource) Int63() int64 {
	s.ncalls++
	return s.r.Int63()
}

func (s *benchmarkSource) Uint64() uint64 {
	s.ncalls++
	return s.r.Uint64()
}

func (s *benchmarkSource) Seed(seed int64) { s.r.Seed(seed) }

func newSource(t testing.TB) (*benchmarkSource, *rand.Rand) {
	seed := time.Now().UnixNano()
	t.Logf("seed: %08x", seed)

	src := &benchmarkSource{
		r: rand.NewSource(seed).(rand.Source64),
	}
	r := rand.New(src)
	return src, r
}

func benchmarkIntn(b *testing.B, n int, intn func(*rand.Rand, rand.Source64, int) int) {
	src, r := newSource(b)

	for i := 0; i < b.N; i++ {
		intn(r, src, n)
	}
	b.ReportMetric(float64(src.ncalls)/float64(b.N), "random_numbers/op")
}

func stdIntn(r *rand.Rand, _ rand.Source64, n int) int      { return r.Intn(n) }
func randstatIntn(_ *rand.Rand, r rand.Source64, n int) int { return randstat.Intn(r, n) }

func BenchmarkIntnStd_16(b *testing.B)  { benchmarkIntn(b, 16, stdIntn) }
func BenchmarkIntnUs_16(b *testing.B)   { benchmarkIntn(b, 16, randstatIntn) }
func BenchmarkIntnStd_123(b *testing.B) { benchmarkIntn(b, 123, stdIntn) }
func BenchmarkIntnUs_123(b *testing.B)  { benchmarkIntn(b, 123, randstatIntn) }
func BenchmarkIntnStd_2k(b *testing.B)  { benchmarkIntn(b, 2000, stdIntn) }
func BenchmarkIntnUs_2k(b *testing.B)   { benchmarkIntn(b, 2000, randstatIntn) }
func BenchmarkIntnStd_1e9(b *testing.B) { benchmarkIntn(b, 1e9, stdIntn) }
func BenchmarkIntnUs_1e9(b *testing.B)  { benchmarkIntn(b, 1e9, randstatIntn) }

func benchmarkInt31n(b *testing.B, n int, intn func(*rand.Rand, rand.Source, int32) int32) {
	src, r := newSource(b)

	for i := 0; i < b.N; i++ {
		intn(r, src, int32(n))
	}
	b.ReportMetric(float64(src.ncalls)/float64(b.N), "random_numbers/op")
}

func stdInt31n(r *rand.Rand, _ rand.Source, n int32) int32      { return r.Int31n(n) }
func randstatInt31n(_ *rand.Rand, r rand.Source, n int32) int32 { return randstat.Int31n(r, n) }

func BenchmarkInt31nStd_16(b *testing.B)  { benchmarkInt31n(b, 16, stdInt31n) }
func BenchmarkInt31nUs_16(b *testing.B)   { benchmarkInt31n(b, 16, randstatInt31n) }
func BenchmarkInt31nStd_123(b *testing.B) { benchmarkInt31n(b, 123, stdInt31n) }
func BenchmarkInt31nUs_123(b *testing.B)  { benchmarkInt31n(b, 123, randstatInt31n) }
func BenchmarkInt31nStd_2k(b *testing.B)  { benchmarkInt31n(b, 2000, stdInt31n) }
func BenchmarkInt31nUs_2k(b *testing.B)   { benchmarkInt31n(b, 2000, randstatInt31n) }
func BenchmarkInt31nStd_1e9(b *testing.B) { benchmarkInt31n(b, 1e9, stdInt31n) }
func BenchmarkInt31nUs_1e9(b *testing.B)  { benchmarkInt31n(b, 1e9, randstatInt31n) }

func benchmarkShuffle(b *testing.B, n int, shuffle func(*rand.Rand, rand.Source64, int)) {
	src, r := newSource(b)

	for i := 0; i < b.N; i++ {
		shuffle(r, src, n)
	}
	b.ReportMetric(float64(src.ncalls)/float64(b.N), "random_numbers/op")
}

func swapNoop(i, j int)                                    {}
func stdShuffle(r *rand.Rand, _ rand.Source64, n int)      { r.Shuffle(n, swapNoop) }
func randstatShuffle(_ *rand.Rand, r rand.Source64, n int) { randstat.Shuffle(r, n, swapNoop) }

func BenchmarkShuffleStd_2k(b *testing.B)   { benchmarkShuffle(b, 2000, stdShuffle) }
func BenchmarkShuffleUs_2k(b *testing.B)    { benchmarkShuffle(b, 2000, randstatShuffle) }
func BenchmarkShuffleStd_20k(b *testing.B)  { benchmarkShuffle(b, 20000, stdShuffle) }
func BenchmarkShuffleUs_20k(b *testing.B)   { benchmarkShuffle(b, 20000, randstatShuffle) }
func BenchmarkShuffleStd_200k(b *testing.B) { benchmarkShuffle(b, 2e5, stdShuffle) }
func BenchmarkShuffleUs_200k(b *testing.B)  { benchmarkShuffle(b, 2e5, randstatShuffle) }
