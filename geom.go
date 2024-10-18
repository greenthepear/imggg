// Modified from go/src/image/geom.go
// With the following copyright notice:
//
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imggg

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"golang.org/x/exp/constraints"
)

// Number interface by combining the exp/contraints' Integer and Float
type Number interface {
	constraints.Integer | constraints.Float
}

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point[V Number] struct {
	X, Y V
}

// Std returns the point as the standard package's image.Point type by a simple cast
// of X and Y to an int.
func (p Point[V]) Std() image.Point {
	return image.Point{int(p.X), int(p.Y)}
}

// String returns a string representation of p like "(3,4)".
func (p Point[V]) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

// Add returns the vector p+q.
func (p Point[V]) Add(q Point[V]) Point[V] {
	return Point[V]{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point[V]) Sub(q Point[V]) Point[V] {
	return Point[V]{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point[V]) Mul(k V) Point[V] {
	return Point[V]{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point[V]) Div(k V) Point[V] {
	return Point[V]{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p Point[V]) In(r Rectangle[V]) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// Mod returns the point q in r such that p.X-q.X is a multiple of r's width
// and p.Y-q.Y is a multiple of r's height.
func (p Point[V]) Mod(r Rectangle[V]) Point[V] {
	w, h := r.Dx(), r.Dy()
	p = p.Sub(r.Min)
	// Since % doesn't work on floats, TODO: avoid cast
	math.Mod(float64(p.X), float64(w))
	if p.X < 0 {
		p.X += w
	}
	math.Mod(float64(p.Y), float64(w))
	if p.Y < 0 {
		p.Y += h
	}
	return p.Add(r.Min)
}

// Eq reports whether p and q are equal.
func (p Point[V]) Eq(q Point[V]) bool {
	return p == q
}

// Pt is shorthand for [Point]{X, Y}.
func Pt[V Number](X, Y V) Point[V] {
	return Point[V]{X, Y}
}

// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
//
// A Rectangle is also an [Image] whose bounds are the rectangle itself. At
// returns color.Opaque for points in the rectangle and color.Transparent
// otherwise.
type Rectangle[V Number] struct {
	Min, Max Point[V]
}

// Std returns the rectangle as the standard package's image.Rectangle type by a simple cast
// of the two points' X and Y to ints.
func (p Rectangle[V]) Std() image.Rectangle {
	return image.Rectangle{
		p.Min.Std(),
		p.Max.Std(),
	}
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle[V]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rectangle[V]) Dx() V {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle[V]) Dy() V {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle[V]) Size() Point[V] {
	return Point[V]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rectangle[V]) Add(p Point[V]) Rectangle[V] {
	return Rectangle[V]{
		Point[V]{r.Min.X + p.X, r.Min.Y + p.Y},
		Point[V]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle[V]) Sub(p Point[V]) Rectangle[V] {
	return Rectangle[V]{
		Point[V]{r.Min.X - p.X, r.Min.Y - p.Y},
		Point[V]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle[V]) Inset(n V) Rectangle[V] {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle[V]) Intersect(s Rectangle[V]) Rectangle[V] {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		return Rectangle[V]{}
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectangle[V]) Union(s Rectangle[V]) Rectangle[V] {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle[V]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle[V]) Eq(s Rectangle[V]) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle[V]) Overlaps(s Rectangle[V]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle[V]) In(s Rectangle[V]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Canon returns the canonical version of r. The returned rectangle has minimum
// and maximum coordinates swapped if necessary so that it is well-formed.
func (r Rectangle[V]) Canon() Rectangle[V] {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

// At implements the [Image] interface.
func (r Rectangle[V]) At(x, y V) color.Color {
	if (Point[V]{x, y}).In(r) {
		return color.Opaque
	}
	return color.Transparent
}

// RGBA64At implements the [RGBA64Image] interface.
func (r Rectangle[V]) RGBA64At(x, y V) color.RGBA64 {
	if (Point[V]{x, y}).In(r) {
		return color.RGBA64{0xffff, 0xffff, 0xffff, 0xffff}
	}
	return color.RGBA64{}
}

// Bounds implements the [Image] interface.
func (r Rectangle[V]) Bounds() Rectangle[V] {
	return r
}

// ColorModel implements the [Image] interface.
func (r Rectangle[V]) ColorModel() color.Model {
	return color.Alpha16Model
}

// Rect is shorthand for [Rectangle]{Pt(x0, y0), [Pt](x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect[V Number](x0, y0, x1, y1 V) Rectangle[V] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[V]{Point[V]{x0, y0}, Point[V]{x1, y1}}
}
