[![Go Reference](https://pkg.go.dev/badge/github.com/greenthepear/imggg.svg)](https://pkg.go.dev/github.com/greenthepear/imggg)
[![Go Report Card](https://goreportcard.com/badge/github.com/greenthepear/imggg)](https://goreportcard.com/report/github.com/greenthepear/imggg)

**imggg** - go's [**im**a**g**e](https://pkg.go.dev/image) **g**eneric **g**eometrics, reimplements [`image.Point`](https://pkg.go.dev/image#Point) and [`image.Rectangle`](https://pkg.go.dev/image#Rectangle) to work with any number type you want: `int`, `int8`, `int16`, `int32`, `int64`, `float32`, `float64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `uintptr`, instead of just `int`.

Features:
- All the same methods and functions as `image`, the code is modified directly from https://go.dev/src/image/geom.go.
- Cast them back to the standard library equivalent with the `Std()` method.
- Get X and Y of as two variables quickly with the `XY()` method.

> [!NOTE]
> Using points and rectangles with floats changes their logic from the intentions of the image library, read more on [The Go Blog about the package](https://go.dev/blog/image).

> [!WARNING]
> Operations on unsigned integers is not tested, just try not to underflow things.

# Usage

```go
package main

import (
    "github.com/greenthepear/imggg"
    "image"
)

func main(){
    // Be implicit
    floatPoint1 := imggg.Pt[float64](0.2,0.1)
    // Let the compiler figure it out
    floatPoint2 := imggg.Pt(1.6,2.2)
    // Use the classic methods
    rec := imggg.Rectangle[float64]{
		floatPoint1.Mul(10),
		floatPoint2.Div(0.1),
	} // (2,1)-(16,22)

    // Work with the standard library
    img := image.NewRGNA(
        rec.Std()
    )
    ...
}
```

# License
Do what you want with this, to be fancy I chose the **0BSD** license.

*Uses Go source code, Copyright 2010 The Go Authors. All rights reserved:* [*LICENSE*](https://cs.opensource.google/go/go/+/refs/tags/go1.23.2:LICENSE)
