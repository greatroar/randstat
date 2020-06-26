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

// +build ignore

package main

import (
	"io"
	"log"
	"os"
	"text/template"
)

func main() {
	var out io.Writer

	switch len(os.Args) {
	case 1:
		out = os.Stdout
	case 3:
		if os.Args[1] != "-o" {
			usage()
		}
		f, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		out = f
	default:
		usage()
	}

	io.WriteString(out, header)

	for _, c := range []struct {
		Bits    int
		Type    string
		Out     string
		SrcBits string
	}{
		{31, "int32", "int", "64"},
		{63, "int64", "int", "64"},
		{31, "int32", "int32", ""},
		{31, "int32", "int64", ""},
		{63, "int64", "int64", "64"},
	} {
		code.Execute(out, c)
	}
}

func usage() { log.Fatalf("usage: %s [-o file]", os.Args[0]) }

const header = `// Code generated by generate_ints.go. DO NOT EDIT.

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

// Algorithm L from Li 1994, Reservoir-Sampling Algorithms of Time
// Complexity O(n(1+log(N/n))), ACM TOMS,
// https://doi.org/10.1145%2F198429.198435,
// specialized for various integer types.

package sampling

import (
	"math"
	"math/rand"

	"github.com/greatroar/randstat"
)
`

var code = template.Must(template.New("code").Parse(`
func ints{{.Bits}}_{{.Out}}(samplesize int, n {{.Type}}, r rand.Source{{.SrcBits}}, buf []{{.Out}}) []{{.Out}} {
	if samplesize == 0 {
		return buf
	}

	for i := 0; i < samplesize; i++ {
		buf = append(buf, {{.Out}}(i))
	}
	sample := buf[len(buf)-samplesize:]

	var (
		w = float64(1)
		i = float64(samplesize)
		k = float64(samplesize)
		N = float64(n)
	)
	for {
		w *= math.Exp(math.Log(random01(r)) / k)
		i += 1 + math.Floor(math.Log(random01(r))/math.Log1p(-w))
		if i >= N {
			break
		}
		j := randstat.Int{{.Bits}}n(r, {{.Type}}(len(sample)))
		sample[j] = {{.Out}}(i)
	}

	return buf
}
`))