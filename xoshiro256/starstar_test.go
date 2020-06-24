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

package xoshiro256_test

import (
	"math/rand"
	"testing"

	"github.com/greatroar/randstat/xoshiro256"

	"github.com/stretchr/testify/assert"
)

var _ rand.Source64 = (*xoshiro256.Source)(nil)

func TestSource(t *testing.T) {
	var s xoshiro256.Source

	for _, c := range []struct {
		seed uint64
		r    [15]uint64
	}{
		// Values generated by the C version, bootstrapped by SplitMix64.
		{seed: 0x00000000,
			r: [...]uint64{
				0x99ec5f36cb75f2b4, 0xbf6e1f784956452a, 0x1a5f849d4933e6e0,
				0x6aa594f1262d2d2c, 0xbba5ad4a1f842e59, 0xffef8375d9ebcaca,
				0x6c160deed2f54c98, 0x8920ad648fc30a3f, 0xdb032c0ba7539731,
				0xeb3a475a3e749a3d, 0x1d42993fa43f2a54, 0x11361bf526a14bb5,
				0x1b4f07a5ab3d8e9c, 0xa7a3257f6986db7f, 0x7efdaa95605dfc9c,
			}},

		{seed: 0x00000001,
			r: [...]uint64{
				0xb3f2af6d0fc710c5, 0x853b559647364cea, 0x92f89756082a4514,
				0x642e1c7bc266a3a7, 0xb27a48e29a233673, 0x24c123126ffda722,
				0x123004ef8df510e6, 0x61954dcc47b1e89d, 0xddfdb48ab9ed4a21,
				0x8d3cdb8c3aa5b1d0, 0xeebd114bd87226d1, 0xf50c3ff1e7d7e8a6,
				0xeeca3115e23bc8f1, 0xab49ed3db4c66435, 0x99953c6c57808dd7,
			}},

		{seed: 0x123456789,
			r: [...]uint64{
				0x57fbcd9b91cb292, 0xfe959e167db74693, 0x1b0ca35f4e529ff6,
				0x58d94faa98fba417, 0x9adef42b023a83ae, 0xcef39808ebe6b956,
				0xcddfeff3f2a26179, 0xf95812addb689ab9, 0xd0b1a2f1db013d6d,
				0x2b938bfa658f761a, 0xcb5d8fe7ace24927, 0x4bda982c991603e3,
				0xf780f90bf3a99b75, 0xfc13f54a3cc9d23e, 0x3008f188e6e7f559,
			}},

		{seed: 0x05f5e0ff,
			r: [...]uint64{
				0xf97b3a1d84a047a6, 0xe6a983a71819be92, 0x12c5f99f0288f12b,
				0x2ea8efe00ac6659, 0x75c04cdc89723635, 0x106d87d418ba2d80,
				0x36f47cfa3c68f0ad, 0x91e34c31fd6ca36b, 0xbf1e77f2248a74a4,
				0x75c150ada67cccf0, 0xcaac6333b57f8354, 0x3ecb5c5882958911,
				0xed3f536de40be757, 0x2ee4dfb7c95d4e2, 0x15102ca8ec0d9373,
			}},
	} {
		s.Seed(int64(c.seed))
		for i, v := range c.r {
			assert.Equal(t, v, s.Uint64(), "seed = 0x%08x, value %d", c.seed, i)
		}
	}
}