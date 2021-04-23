package sphericalmercator

import (
	"testing"
)

var sm SphericalMercator

var MAXEXTENTMERC = []float64{-20037508.342789244, -20037508.342789244, 20037508.342789244, 20037508.342789244}
var MAXEXTENTWGS84 = []float64{-180, -85.0511287798066, 180, 85.0511287798066}

func init() {
	sm = New(nil)
}

func TestBBox(t *testing.T) {
	want1 := []float64{-180, -85.05112877980659, 180, 85.0511287798066}
	got1 := sm.BBox(0, 0, 0, true, "WGS84")
	if !sliceEqual(got1, want1) {
		t.Errorf("got %v; want %v\n", got1, want1)
	}

	want2 := []float64{-180, -85.05112877980659, 0, 0}
	got2 := sm.BBox(0, 0, 1, true, "WGS84")
	if !sliceEqual(got2, want2) {
		t.Errorf("got %v; want %v\n", got2, want2)
	}
}

func TestXYZ(t *testing.T) {
	want1 := Bounds{MinX: 0, MinY: 0, MaxX: 0, MaxY: 0}
	got1 := sm.XYZ([]float64{-180, -85.05112877980659, 180, 85.0511287798066}, 0, true, "WGS84")
	if got1 != want1 {
		t.Errorf("got %v; want %v\n", got1, want1)
	}

	want2 := Bounds{MinX: 0, MinY: 0, MaxX: 0, MaxY: 0}
	got2 := sm.XYZ([]float64{-180, -85.05112877980659, 0, 0}, 1, true, "WGS84")
	if got2 != want2 {
		t.Errorf("got %v; want %v\n", got2, want2)
	}
}

func TestXYZBroken(t *testing.T) {
	extent := []float64{-0.087891, 40.95703, 0.087891, 41.044916}
	got := sm.XYZ(extent, 3, true, "WGS84")
	if !(got.MinX <= got.MaxX) {
		t.Errorf("x: %f > %f for %v\n", got.MinX, got.MaxX, extent)
	}
	if !(got.MinY <= got.MaxY) {
		t.Errorf("x: %f > %f for %v\n", got.MinY, got.MaxY, extent)
	}
}

func TestXYZNegative(t *testing.T) {
	extent := []float64{-112.5, 85.0511, -112.5, 85.0511}
	got := sm.XYZ(extent, 0, false, "")
	if got.MinY != 0 {
		t.Errorf("min y should be zero\n")
	}
}

// ToDo: TestXYZFuzz

func TestConvert(t *testing.T) {
	want1 := MAXEXTENTMERC
	got1 := sm.Convert(MAXEXTENTWGS84, "900913")
	if !sliceEqual(got1, want1) {
		t.Errorf("got %v; want %v\n", got1, want1)
	}

	want2 := MAXEXTENTWGS84
	got2 := sm.Convert(MAXEXTENTMERC, "WGS84")
	if !sliceEqual(got2, want2) {
		t.Errorf("got %v; want %v\n", got2, want2)
	}
}

func TestExtents(t *testing.T) {
	want1 := MAXEXTENTMERC
	got1 := sm.Convert([]float64{-240, -90, 240, 90}, "900913")
	if !sliceEqual(got1, want1) {
		t.Errorf("got %v; want %v\n", got1, want1)
	}

	want2 := Bounds{MinX: 0, MinY: 0, MaxX: 15, MaxY: 15}
	got2 := sm.XYZ([]float64{-240, -90, 240, 90}, 4, true, "WGS84")
	if got2 != want2 {
		t.Errorf("got %v; want %v\n", got2, want2)
	}
}

func TestLL(t *testing.T) {
	want1 := []float64{-179.45068359375, 85.00351401304403}
	got1 := sm.LL([]float64{200, 200}, 9)
	if !sliceEqual(got1, want1) {
		t.Errorf("got %v; want %v\n", got1, want1)
	}

	want2 := []float64{-179.3034449476476, 84.99067388699072}
	got2 := sm.LL([]float64{200, 200}, 8.6574)
	if !sliceEqual(got2, want2) {
		t.Errorf("got %v; want %v\n", got2, want2)
	}
}

func TestPx(t *testing.T) {
	want1 := []float64{364, 215}
	got1 := sm.Px([]float64{-179, 85}, 9)
	if !sliceEqual(got1, want1) {
		t.Errorf("got %v; want %v\n", got1, want1)
	}

	want2 := []float64{287.12734093961626, 169.30444219392666}
	got2 := sm.Px([]float64{-179, 85}, 8.6574)
	if !sliceEqual(got2, want2) {
		t.Errorf("got %v; want %v\n", got2, want2)
	}
}

// ToDo: TestHighPrecisionFloat

func sliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
