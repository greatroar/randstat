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

// Package randstat provides replacements for math/rand functions and methods.
//
// The functions in this package are designed to be faster than their
// counterparts in math/rand, and consume (in expectation) fewer random numbers
// from their Source. Because they use different algorithms, they produce
// different sequences of random numbers.
//
// The subpackages provide various random number generators and sampling
// algorithms.
package randstat
