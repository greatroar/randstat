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

// Derived from http://prng.di.unimi.it/xoshiro256starstar.c:
//
// Written in 2018 by David Blackman and Sebastiano Vigna (vigna@acm.org)
//
// To the extent possible under law, the author has dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
//
// See <http://creativecommons.org/publicdomain/zero/1.0/>.

// Package xoshiro256 implements the xoshiro256** random number generator.
package xoshiro256

import (
	"encoding/binary"
	"errors"
	"math/bits"

	"github.com/greatroar/randstat/splitmix64"
)

// A Source is a xoshiro256** 1.0 random number generator.
//
// A Source must be seeded (with the Seed method) before use.
type Source struct{ s [4]uint64 }

// New returns a Source initialized with the given seed.
func New(seed uint64) *Source {
	s := &Source{}
	s.Seed(int64(seed))
	return s
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (s *Source) Int63() int64 { return int64(s.Uint64() >> 1) }

// Jump advances the Source by 2^128 positions.
func (s *Source) Jump() {
	jump := [...]uint64{
		0x180ec6d33cfd0aba, 0xd5a61266f0c9392c,
		0xa9582618e03fc9aa, 0x39abdc4529b1661c}

	var s0, s1, s2, s3 uint64

	for i := 0; i < len(jump); i++ {
		for b := byte(0); b < 64; b++ {
			if jump[i]&(1<<b) != 0 {
				s0 ^= s.s[0]
				s1 ^= s.s[1]
				s2 ^= s.s[2]
				s3 ^= s.s[3]
			}
			s.Uint64()
		}
	}

	s.s = [4]uint64{s0, s1, s2, s3}
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
//
// It uses a SplitMix64 generator to turn seed into four non-zero
// pseudo-random numbers.
func (s *Source) Seed(seed int64) {
	sm := splitmix64.Source(seed)

retry:
	s0 := sm.Uint64()
	s1 := sm.Uint64()
	s2 := sm.Uint64()
	s3 := sm.Uint64()

	if s0|s1|s2|s3 == 0 {
		goto retry
	}

	s.s = [4]uint64{s0, s1, s2, s3}
}

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (s *Source) Uint64() uint64 {
	st := &s.s
	r := bits.RotateLeft64(5*st[1], 7) * 9

	t := st[1] << 17
	st[2] ^= st[0]
	st[3] ^= st[1]
	st[1] ^= st[2]
	st[0] ^= st[3]
	st[2] ^= t

	st[3] = bits.RotateLeft64(st[3], 45)

	return r
}

const (
	header      = "xoshiro256**1.0\x00"
	marshalSize = len(header) + 4*8
)

// MarshalBinary encodes s in a binary format for serialization.
//
// The format starts with a 16-byte header, followed by the four 64-bit
// integers of state in little-endian format.
//
// The returned error is always nil.
func (s *Source) MarshalBinary() (data []byte, err error) {
	data = make([]byte, marshalSize)
	copy(data, header)
	for i, x := range s.s {
		binary.LittleEndian.PutUint64(data[len(header)+8*i:], x)
	}
	return
}

// UnmarshalBinary decodes s from the binary format used by MarshalBinary.
func (s *Source) UnmarshalBinary(data []byte) error {
	switch {
	case len(data) != marshalSize:
		return errors.New("unmarshal xoshiro256: incorrect data length")
	case string(data[:len(header)]) != header:
		return errors.New("unmarshal xoshiro256: incorrect header")
	}

	data = data[len(header):]
	for i := range s.s {
		s.s[i] = binary.LittleEndian.Uint64(data)
		data = data[8:]
	}
	return nil
}
