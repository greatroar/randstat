// +build go1.12

package randstat

import "math/bits"

func mul64(a, b uint64) (hi, lo uint64) {
	return bits.Mul64(a, b)
}
