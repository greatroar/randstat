package sampling_test

import (
	"math/rand"

	"github.com/greatroar/randstat/sampling"
)

func ExampleVaropt_Show_reuseMemory() {
	// The return value from Show can be used to reuse allocated memory.

	type item struct {
		letter rune
		index  int
		weight float64
	}

	sample := sampling.NewVaropt(6, nil)
	var x *item

	for i, l := range []rune("abcdefghijklmnopqrstuvwxyz") {
		if x == nil {
			x = new(item)
		}
		*x = item{letter: l, index: i, weight: rand.Float64()}

		reject := sample.Show(x, x.weight)

		// A rejected item is no longer in the sample and may be recycled.
		x, _ = reject.(*item)
	}
	//Output:
}
