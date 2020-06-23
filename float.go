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

import "math/rand"

// In Go 1.13, we could write 0x1p-53.
const machineEpsilon = 1. / (1. << 53)

// Float64 returns a pseudo-random number in the interval [0,1),
// given a uniformly random integer.
//
// Unlike the rand.Rand.Float64 method, this function consumes exactly
// one random integer from the source r.
func Float64(r rand.Source) float64 {
	// Vigna's recipe from http://prng.di.unimi.it/,
	// adapted to non-negative int63 values.
	return float64(r.Int63()>>10) * machineEpsilon
}
