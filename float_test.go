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
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/greatroar/randstat"
	"github.com/greatroar/randstat/internal/source"

	"github.com/stretchr/testify/assert"
)

func TestFloat64(t *testing.T) {
	t.Parallel()

	x := randstat.Float64(source.Constant(0))
	assert.Equal(t, float64(0), x)

	// Smallest float larger than the IEEE machine epsilon.
	Δ := math.Nextafter(1./(1<<53), 1)

	x = randstat.Float64(source.Constant(1<<63 - 1))
	assert.Less(t, x, float64(1))
	assert.InDelta(t, 1, x, Δ)

	r := rand.NewSource(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		u := r.Int63()
		x := randstat.Float64(source.Constant(u))
		y := randstat.Float64(source.Constant(1 + u))
		assert.InDelta(t, x, y, Δ)
	}
}

func BenchmarkFloat64Std(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		r.Float64()
	}
}

func BenchmarkFloat64Us(b *testing.B) {
	r := rand.NewSource(time.Now().UnixNano())

	b.ResetTimer()
	b.SetBytes(8)

	for i := 0; i < b.N; i++ {
		randstat.Float64(r)
	}
}
