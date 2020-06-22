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

	"github.com/greatroar/randstat"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShuffle(t *testing.T) {
	t.Parallel()

	const n = 9999

	var ncalls int
	swap := func(i, j int) {
		ncalls++
		require.Less(t, i, n)
		require.Less(t, j, i)
	}

	require.Panics(t, func() {
		randstat.Shuffle(rand.NewSource(1).(rand.Source64), -1, swap)
	})
	require.NotPanics(t, func() {
		randstat.Shuffle(rand.NewSource(1).(rand.Source64), 0, swap)
	})

	randstat.Shuffle(rand.NewSource(1).(rand.Source64), n, swap)
	require.Less(t, ncalls, n)
}

func TestShuffleAllPermutations(t *testing.T) {
	t.Parallel()

	var (
		a     = []byte{1, 2, 3, 4, 5, 6, 7}
		perms = make(map[string]struct{})
		r     = rand.NewSource(126).(rand.Source64)
	)

	for len(perms) < 5040 {
		randstat.Shuffle(r, len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
		perms[string(a)] = struct{}{}
	}
}

func TestInt31n(t *testing.T) {
	t.Parallel()

	r := rand.NewSource(0x6983623)

	for _, max := range []int32{1, 163, 1698687, 1 << 20, 1<<31 - 1} {
		for i := 0; i < 20000; i++ {
			require.Less(t, randstat.Int31n(r, max), max)
		}
	}
}

func TestInt63n(t *testing.T) {
	t.Parallel()

	r := rand.NewSource(0x6983641).(rand.Source64)

	for _, max := range []int64{1, 2, 198687, 1 << 32, 1 << 42, 1<<60 + 1} {
		for i := 0; i < 20000; i++ {
			require.Less(t, randstat.Int63n(r, max), max)
		}
	}
}

// Very quick statistical check.
func TestInt31nSmallStats(t *testing.T) {
	t.Parallel()

	const (
		N      = 20
		rounds = 1000
	)

	r := rand.NewSource(0xae691)
	freq := make([]int, N)

	for i := 0; i < rounds*N; i++ {
		freq[randstat.Int31n(r, N)]++
	}
	for _, f := range freq {
		assert.InEpsilon(t, rounds, f, .08)
	}
}

// Very quick statistical check.
func TestInt63nSmallStats(t *testing.T) {
	t.Parallel()

	const (
		N      = 20
		rounds = 1000
	)

	r := rand.NewSource(0x12112).(rand.Source64)
	freq := make([]int, N)

	for i := 0; i < rounds*N; i++ {
		freq[randstat.Int63n(r, N)]++
	}
	for _, f := range freq {
		assert.InEpsilon(t, rounds, f, .08)
	}
}
