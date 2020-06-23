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

package source

import "math/rand"

type constant struct{ C int64 }

func Constant(seed int64) rand.Source { return &constant{C: seed} }
func (c *constant) Int63() int64      { return int64(c.C) }
func (c *constant) Seed(seed int64)   { c.C = seed }

type std struct{}

var Std *std

func (*std) Int63() int64    { return rand.Int63() }
func (*std) Uint64() uint64  { return rand.Uint64() }
func (*std) Seed(seed int64) { rand.Seed(seed) }
