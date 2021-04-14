package sphericalmercator

import (
	"log"
	"math"
)

const (
	epsln = 1.0e-10
	d2r   = math.Pi / 180.0
	r2d   = 180.0 / math.Pi
	// 900913 properties
	a         = 6378137.0
	maxExtent = 20037508.342789244
)

type Cache map[int]map[string][]float64

var (
	cache = Cache{}
)

type Options struct {
	Size int
}

type SphericalMercator struct {
	size  int
	sizef float64
	bc    []float64
	cc    []float64
	zc    []float64
	ac    []float64
}

func New(opts *Options) SphericalMercator {

	sm := SphericalMercator{}

	if opts == nil {
		sm.size = 256
	} else {
		sm.size = opts.Size
	}
	sm.sizef = float64(sm.size)

	if _, ok := cache[sm.size]; !ok {
		size := sm.sizef
		c := map[string][]float64{}
		c["bc"] = make([]float64, 30)
		c["cc"] = make([]float64, 30)
		c["zc"] = make([]float64, 30)
		c["ac"] = make([]float64, 30)

		for d := 0; d < 30; d++ {
			c["bc"][d] = size / 360.0
			c["cc"][d] = size / (2.0 * math.Pi)
			c["zc"][d] = size / 2.0
			c["ac"][d] = size
			size *= 2
		}
		cache[sm.size] = c

		sm.bc = cache[sm.size]["bc"]
		sm.cc = cache[sm.size]["cc"]
		sm.zc = cache[sm.size]["zc"]
		sm.ac = cache[sm.size]["ac"]
	}

	return sm
}

// Px converts an lon/lat pair to a x/y pair.
func (sm SphericalMercator) Px(ll []float64, zoom interface{}) []float64 {
	switch z := zoom.(type) {
	case float64:
		size := sm.sizef * math.Pow(2, z)
		d := size / 2.0
		bc := size / 360.0
		cc := size / (2.0 * math.Pi)
		ac := size
		f := math.Min(math.Max(math.Sin(d2r*ll[1]), -0.9999), 0.9999)
		x := d + ll[0]*bc
		y := d + 0.5*math.Log((1+f)/(1-f))*-cc
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Logical_AND_assignment
		if x > ac {
			x = ac
		}
		if y > ac {
			y = ac
		}
		// (x>ac) && (x = ac)
		// (y>ac) && (y = ac)
		// (x>ac) && (x = 0)
		// (x>ac) && (x = 0)
		return []float64{x, y}
	case int:
		d := sm.zc[z]
		f := math.Min(math.Max(math.Sin(d2r*ll[1]), -0.9999), 0.9999)
		x := math.Round(d + ll[0]*sm.bc[z])
		y := math.Round(d + 0.5*math.Log((1+f)/(1-f))*(-sm.cc[z]))
		if x > sm.ac[z] {
			x = sm.ac[z]
		}
		if y > sm.ac[z] {
			y = sm.ac[z]
		}
		// (x > this.Ac[zoom]) && (x = this.Ac[zoom]);
		// (y > this.Ac[zoom]) && (y = this.Ac[zoom]);
		//(x < 0) && (x = 0);
		//(y < 0) && (y = 0);
		return []float64{x, y}
	default:
		log.Printf("zoom (%T) not supported", zoom)
		return nil
	}
}

// LL converts an x/y pair to an lon/lat pair.
func (sm SphericalMercator) LL(px []float64, zoom interface{}) []float64 {
	switch z := zoom.(type) {
	case float64:
		size := sm.sizef * math.Pow(2, z)
		bc := size / 360.0
		cc := size / (2.0 * math.Pi)
		zc := size / 2.0
		g := (px[1] - zc) / -cc
		lon := (px[0] - zc) / bc
		lat := r2d * (2.0*math.Atan(math.Exp(g)) - 0.5*math.Pi)
		return []float64{lon, lat}
	case int:
		g := (px[1] - sm.zc[z]) / (-sm.cc[z])
		lon := (px[0] - sm.zc[z]) / sm.bc[z]
		lat := r2d * (2.0*math.Atan(math.Exp(g)) - 0.5*math.Pi)
		return []float64{lon, lat}
	default:
		log.Printf("zoom (%T) not supported", zoom)
		return nil
	}
}

// BBox converts tile xyz value to bbox of the form `[w, s, e, n]`.
func (sm SphericalMercator) BBox(x, y float64, zoom interface{}, tmsStyle bool, srs string) []float64 {

	// Convert xyz into bbox with srs WGS84
	if tmsStyle {
		y = (math.Pow(2, zz(zoom)) - 1) - y
	}

	ll := []float64{x * sm.sizef, (y + 1) * sm.sizef} // lower left
	ur := []float64{(x + 1) * sm.sizef, y * sm.sizef} // upper right

	var bbox []float64
	bbox = append(bbox, sm.LL(ll, zoom)...)
	bbox = append(bbox, sm.LL(ur, zoom)...)

	// If web mercator requested reproject to 900913.
	if srs == "900913" {
		return sm.Convert(bbox, "900913")
	}

	return bbox
}

type Bounds struct {
	MinX, MinY float64
	MaxX, MaxY float64
}

// XYZ converts bbox to xyx bounds.
func (sm SphericalMercator) XYZ(bbox []float64, zoom interface{}, tmsStyle bool, srs string) Bounds {
	// If web mercator provided reproject to WGS84.
	if srs == "900913" {
		bbox = sm.Convert(bbox, "WGS84")
	}

	ll := []float64{bbox[0], bbox[1]} // lower left
	ur := []float64{bbox[2], bbox[3]} // upper right
	pxLL := sm.Px(ll, zoom)
	pxUR := sm.Px(ur, zoom)
	// Y = 0 for XYZ is the top hence minY uses px_ur[1].
	x := []float64{math.Floor(pxLL[0] / sm.sizef), math.Floor((pxUR[0] - 1) / sm.sizef)}
	y := []float64{math.Floor(pxUR[1] / sm.sizef), math.Floor((pxLL[1] - 1) / sm.sizef)}
	bounds := Bounds{
		MinX: minZero(x),
		MinY: minZero(y),
		MaxX: max(x),
		MaxY: max(y),
	}
	if tmsStyle {
		tms := Bounds{
			MinY: (math.Pow(2, zz(zoom)) - 1) - bounds.MaxY,
			MaxY: (math.Pow(2, zz(zoom)) - 1) - bounds.MinY,
		}
		bounds.MinY = tms.MinY
		bounds.MaxY = tms.MaxY
	}
	return bounds
}

// Convert converts the projection of a given bbox.
func (sm SphericalMercator) Convert(bbox []float64, to string) []float64 {
	if to == "900913" {
		return append(sm.Forward(bbox[0:2]), sm.Forward(bbox[2:4])...)
	}
	return append(sm.Inverse(bbox[0:2]), sm.Inverse(bbox[2:4])...)
}

// Forward converts lon/lat values to 900913 x/y.
func (sm SphericalMercator) Forward(ll []float64) []float64 {
	xy := []float64{
		a * ll[0] * d2r,
		a * math.Log(math.Tan((math.Pi*0.25)+(0.5*ll[1]*d2r))),
	}
	// if xy value is beyond maxextent (e.g. poles), return maxextent.
	if xy[0] > maxExtent {
		xy[0] = maxExtent
	}
	if xy[0] < -maxExtent {
		xy[0] = -maxExtent
	}
	if xy[1] > maxExtent {
		xy[1] = maxExtent
	}
	if xy[1] < -maxExtent {
		xy[1] = -maxExtent
	}
	return xy
}

// Inverse converts 900913 x/y values to lon/lat.
func (sm SphericalMercator) Inverse(xy []float64) []float64 {
	return []float64{
		xy[0] * r2d / a,
		((math.Pi * 0.5) - 2.0*math.Atan(math.Exp(-xy[1]/a))) * r2d,
	}
}

func minZero(a []float64) float64 {
	m := a[0]
	for _, v := range a {
		if v < m {
			m = v
		}
	}
	if m < 0 {
		return 0
	}
	return m
}

func max(a []float64) float64 {
	m := a[0]
	for _, v := range a {
		if v > m {
			m = v
		}
	}
	return m
}

func zz(zoom interface{}) float64 {
	switch z := zoom.(type) {
	case float64:
		return z
	case int:
		return float64(z)
	default:
		log.Printf("zoom (%T) not supported", zoom)
		return 0
	}
}
