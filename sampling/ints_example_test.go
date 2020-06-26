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
	// 209976569432343488
	// 512273708663700480
	// 204888252163057856
	// 311815901727164224
	// 943870718497463296
	// 188287930289750272
	// 629183818546644096
	// 208260206501070304
	// 937471898537404672
	// 639599965878877696
}
