package sampling

import (
	"container/heap"
	"math/rand"

	"github.com/greatroar/randstat"
	"github.com/greatroar/randstat/internal/sums"
)

// A Varopt is a weighted reservoir sampler. Items are included in its sample
// with probability proportional to their weight.
//
// Varopt implements the algorithm of https://arxiv.org/pdf/0803.0473.pdf.
type Varopt struct {
	// Reservoir, divided according to weight above or below the threshold.
	large minWeight
	small []item

	smallbuf []item

	r         rand.Source64
	size      int
	threshold float64
}

// NewVaropt constructs a Varopt sampler.
//
// Random numbers are taken from r, or from an internal generator
// bootstrapped from math.rand's global generator if r is nil.
func NewVaropt(samplesize int, r rand.Source64) *Varopt {
	if samplesize < 0 {
		panic("negative sample size")
	}

	r = maybeXoshiro(r)

	reservoir := make([]item, 3+3*samplesize)
	small := reservoir[:1+samplesize]
	large := reservoir[1+samplesize : 2+2*samplesize]
	smallbuf := reservoir[2+2*samplesize:]

	return &Varopt{
		large: large[: 0 : 1+samplesize],
		small: small[: 0 : 1+samplesize],

		smallbuf: smallbuf[: 0 : 1+samplesize],

		r:    r,
		size: samplesize,
	}
}

type item struct {
	v interface{}
	w float64
}

// Show presents x to v as a candidate for inclusion in its random sample.
//
// An item with zero weight is always rejected. A negative weight causes Show
// to panic.
//
// If x is accepted into the sample, reject is set to the item evicted to make
// space for it, if any. For the first samplesize items, reject will be nil.
func (v *Varopt) Show(x interface{}, w float64) (reject interface{}) {
	switch {
	case w == 0:
		return x

	case w < 0:
		panic("negative weight")

	case v.Len() < v.size:
		v.large = append(v.large, item{x, w})
		if len(v.large) == v.size {
			heap.Init(&v.large)
		}
		return nil
	}

	//var wsum sums.Neumaier
	var wsum sums.Kahan
	wsum.Set(v.threshold * float64(len(v.small)))

	small := v.smallbuf[:0]

	if w > v.threshold {
		// heap.Push(&v.large, item{x, w})
		v.large = append(v.large, item{x, w})
		heap.Fix(&v.large, len(v.large)-1)
	} else {
		small = append(small, item{x, w})
		wsum.Add(w)
	}

	for len(v.large) > 0 &&
		wsum.Value() >= float64(len(v.small)+len(small)-1)*v.large[0].w {

		// small = append(small, heap.Pop(&v.large).(item))
		var lmin item
		lmin, v.large = remove(v.large, 0)
		small = append(small, lmin)
		heap.Fix(&v.large, 0)
		wsum.Add(lmin.w)
	}

	t := wsum.Value() / float64(len(v.small)+len(small)-1)
	var r sums.Kahan
	r.Set(random01(v.r))

	j := 0
	for j < len(small) && r.Value() >= 0 {
		r.Add(-1)
		r.Add(small[j].w / t)
		j++
	}

	var evict item
	if r.Value() < 0 {
		evict, small = remove(small, j-1)
	} else {
		j = randstat.Intn(v.r, len(v.small))
		evict, v.small = remove(v.small, j)
	}

	v.small = append(v.small, small...)
	for i := range small {
		small[i].v = nil // Allow garbage collection.
	}
	v.smallbuf = small
	v.threshold = t

	return evict.v
}

// Item returns the item at index i in the current sample.
//
// The index i must be at least zero and less than s.Len().
// Items occur in the sample in random order.
func (v *Varopt) Item(i int) interface{} {
	if i < len(v.large) {
		return v.large[i].v
	}
	return v.small[i-len(v.large)].v
}

// Len returns the number of items currently in the sample.
//
// The number of items is the minimum of the desired sample size
// and the number of items Shown with positive weight.
func (v *Varopt) Len() int { return len(v.large) + len(v.small) }

func remove(a []item, i int) (item, []item) {
	x := a[i]
	n := len(a) - 1
	a[i] = a[n]
	a[n].v = nil // Allow garbage collection.
	return x, a[:n]
}

// Min-priority queue keyed on weight.
type minWeight []item

func (h *minWeight) Len() int           { return len(*h) }
func (h *minWeight) Less(i, j int) bool { return (*h)[i].w < (*h)[j].w }
func (*minWeight) Pop() interface{}     { panic("use heap.Fix") }
func (*minWeight) Push(interface{})     { panic("use heap.Fix") }
func (h *minWeight) Swap(i, j int)      { a := *h; a[i], a[j] = a[j], a[i] }
