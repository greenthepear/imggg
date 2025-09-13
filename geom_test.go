package imggg

import (
	"image"
	"testing"
)

// TODO: Tests maybe lol
func TestGeom(t *testing.T) {
	// Readme example test
	floatPoint1 := Pt(0.2, 0.1)
	floatPoint2 := Pt(1.6, 2.2)
	rec := Rectangle[float64]{
		floatPoint1.Mul(10),
		floatPoint2.Div(0.1),
	}
	img := image.NewRGBA(
		rec.Std(),
	)
	_ = img

	if Pt(2, 2).Dot(Pt(3, 3)) != 12 {
		t.Errorf("wrong dot product")
	}
}
