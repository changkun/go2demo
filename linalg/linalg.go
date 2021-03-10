// Copyright 2020 Changkun Ou. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linalg

import "math"

// Numeric is a contract that matches any numeric type.
// It would likely be in a contracts package in the standard library.
type Numeric interface {
	type int, int8, int16, int32, int64, uint, uint8,
		uint16, uint32, uint64, uintptr, float32,
		float64, complex64, complex128
}

// NumericAbs matches numeric types with an Abs method.
type NumericAbs[T any] interface {
	Numeric
	Abs() T
}

// OrderedNumeric matches numeric types that support the < operator.
type OrderedNumeric interface {
	type int, int8, int16, int32, int64, uint, uint8,
		uint16, uint32, uint64, uintptr, float32, float64
}

// Complex matches the two complex types, which do not have a < operator.
type Complex interface {
	type complex64, complex128
}

func DotProduct[T Numeric](s1, s2 []T) T {
	if len(s1) != len(s2) {
		panic("DotProduct: slices of unequal length")
	}
	var r T
	for i := range s1 {
		r += s1[i] * s2[i]
	}
	return r
}

// AbsDifference computes the absolute value of the difference of
// a and b, where the absolute value is determined by the Abs method.
func AbsDifference[T NumericAbs](a, b T) T {
	d := a - b
	return d.Abs()
}

// OrderedAbs is a helper type that defines an Abs method for
// ordered numeric types.
type OrderedAbs[T OrderedNumeric] T

func (a OrderedAbs[T]) Abs() OrderedAbs[T] {
	if a < 0 {
		return -a
	}
	return a
}

// ComplexAbs is a helper type that defines an Abs method for
// complex types.
type ComplexAbs[T Complex] T

func (a ComplexAbs[T]) Abs() ComplexAbs[T] {
	r := float64(real(a))
	i := float64(imag(a))
	d := math.Sqrt(r * r + i * i)
	return ComplexAbs[T](complex(d, 0))
}

func OrderedAbsDifference[T OrderedNumeric](a, b T) T {
	return T(AbsDifference(OrderedAbs[T](a), OrderedAbs[T](b)))
}

func ComplexAbsDifference[T Complex](a, b T) T {
	return T(AbsDifference(ComplexAbs[T](a), ComplexAbs[T](b)))
}