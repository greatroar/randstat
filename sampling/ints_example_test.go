package sampling_test

import (
	"fmt"
	"math/rand"

	"github.com/greatroar/randstat/sampling"
)

func ExampleInts64() {
	// Ints64 can sample efficiently from very large ranges.
	// The final argument may be nil.
	r := rand.NewSource(42).(rand.Source64)
	sample := sampling.Ints64(10, 1e18, r, nil)

	for _, x := range sample {
		fmt.Println(x)
	}

	// Output:
	// 209976569432343232
	// 512273708663699840
	// 204888252163057600
	// 311815901727163776
	// 943870718497462016
	// 188287930289750048
	// 629183818546643200
	// 208260206501070048
	// 937471898537403392
	// 639599965878876800
}
