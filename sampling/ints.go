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

//go:generate go run generate_ints.go -o zints.go

package sampling

import (
	"math/rand"

	"github.com/greatroar/randstat/xoshiro256"
)

const maxint32 = 1<<31 - 1

// Ints appends to buf a simple random sample of the integers [0,n)
// and returns the resulting slice. The sample is not sorted.
//
// Random numbers are taken from r, or from an internal generator bootstrapped
// from math.rand's global generator if r is nil.
//
// If samplesize > n, the sample will be of size n instead.
//
// The time complexity of this function is O(s(1+log(n/s))).
func Ints(samplesize, n int, r rand.Source64, buf []int) []int {
	r = maybeXoshiro(r)
	if samplesize > n {
		samplesize = n
	}

	if n <= maxint32 {
		return ints31_int(samplesize, int32(n), r, buf)
	} else {
		return ints63_int(samplesize, int64(n), r, buf)
	}
}

// Ints32 appends to buf a simple random sample of the integers [0,n)
// and returns the resulting slice. The sample is not sorted.
//
// Random numbers are taken from r, or from an internal generator bootstrapped
// from math.rand's global generator if r is nil.
//
// If samplesize > n, the sample will be of size n instead.
//
// The time complexity of this function is O(s(1+log(n/s))).
func Ints32(samplesize int, n int32, r rand.Source, buf []int32) []int32 {
	if r == nil {
		r = xoshiro256.New(rand.Uint64())
	}
	if samplesize > int(n) {
		samplesize = int(n)
	}
	return ints31_int32(samplesize, n, r, buf)
}

// Ints64 appends to buf a simple random sample of the integers [0,n)
// and returns the resulting slice. The sample is not sorted.
//
// Random numbers are taken from r, or from an internal generator bootstrapped
// from math.rand's global generator if r is nil.
//
// If samplesize > n, the sample will be of size n instead.
//
// The time complexity of this function is O(s(1+log(n/s))).
func Ints64(samplesize int, n int64, r rand.Source64, buf []int64) []int64 {
	r = maybeXoshiro(r)
	if int64(samplesize) > n {
		samplesize = int(n)
	}

	if n <= maxint32 {
		return ints31_int64(samplesize, int32(n), r, buf)
	} else {
		return ints63_int64(samplesize, n, r, buf)
	}
}
