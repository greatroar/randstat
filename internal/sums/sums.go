// Package sums implements summation algorithms.
//
// The algorithms in this package put reproducibility before precision,
// so they may round too often to prevent differences between platforms
// and compiler versions.
package sums

// The float64 conversions force rounding to occur.
// See Go spec, §Floating-point operators.

import "math"

func abs(x float64) float64 { return math.Abs(x) }

// Classic Kahan summation.
type Kahan struct{ sum, err float64 }

func (s *Kahan) Add(x float64) {
	x = float64(x - s.err)
	sum := float64(s.sum + x)
	s.err = float64(float64(sum-s.sum) - x)
	s.sum = sum
}

func (s *Kahan) Set(x float64) { s.sum, s.err = x, 0 }

func (s *Kahan) Value() float64 { return float64(s.sum + s.err) }

// Neumaier's improved Kahan–Babuška summation,
// https://www.mat.univie.ac.at/~neum/scan/01.pdf.
type Neumaier struct{ sum, err float64 }

func (s *Neumaier) Add(x float64) {
	t := float64(s.sum + x)
	if abs(s.sum) >= abs(x) {
		s.err += float64(float64(s.sum-t) + x)
	} else {
		s.err += float64(float64(x-t) + s.sum)
	}
	s.sum = t
}

func (s *Neumaier) Set(x float64) { s.sum, s.err = x, 0 }

func (s *Neumaier) Value() float64 { return float64(s.sum + s.err) }
