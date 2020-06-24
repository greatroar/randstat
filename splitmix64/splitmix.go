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

// Derived from http://vigna.di.unimi.it/xorshift/splitmix64.c:
//
// Written in 2015 by Sebastiano Vigna (vigna@acm.org)
//
// To the extent possible under law, the author has dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
//
// See <http://creativecommons.org/publicdomain/zero/1.0/>.

// Package splitmix64 implements the SplitMix64 random number generator.
package splitmix64

import "sync/atomic"

// A Source is a SplitMix64 random number generator.
type Source uint64

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (s *Source) Int63() int64 { return int64(s.Uint64() >> 1) }

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (s *Source) Seed(seed int64) { *s = Source(seed) }

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (s *Source) Uint64() uint64 {
	*s += 0x9e3779b97f4a7c15
	z := uint64(*s)
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// An AtomicSource is a concurrency-safe SplitMix64 random number generator.
//
// All methods on an AtomicSource may safely be called concurrently by
// multiple goroutines.
type AtomicSource uint64

// Int63 returns a 63-bit random number.
func (s *AtomicSource) Int63() int64 { return int64(s.Uint64() >> 1) }

// Seed seeds the generator.
func (s *AtomicSource) Seed(seed int64) {
	atomic.StoreUint64((*uint64)(s), uint64(seed))
}

// Uint64Atomic returns a 64-bit random number.
// It may be safely called by multiple goroutines concurrently.
func (s *AtomicSource) Uint64() uint64 {
	z := atomic.AddUint64((*uint64)(s), 0x9e3779b97f4a7c15)
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}
