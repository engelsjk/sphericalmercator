# sphericalmercator

Projection math for conversions between mercator meters, screen pixels and lat/lon pairs. A Go port of [mapbox/sphericalmercator](https://github.com/mapbox/sphericalmercator).

## Installation

```bash
go get github.com/engelsjk/sphericalmercator
```

## API

```Go
import "github.com/engelsjk/sphericalmercator"

sm := sphericalmercator.New(&sphericalmercator.Options{Size: 512}) 

// Px converts an lon/lat pair to a x/y pair.
sm.Px(ll []float64, zoom interface{}) []float64

// LL converts an x/y pair to an lon/lat pair.
sm.LL(px []float64, zoom interface{}) []float64

// BBox converts tile xyz value to bbox of the form `[w, s, e, n]`.
sm.BBox(x, y float64, zoom interface{}, tmsStyle bool, srs string) []float64

// XYZ converts bbox to xyx bounds.
sm.XYZ(bbox []float64, zoom interface{}, tmsStyle bool, srs string) Bounds

// Convert converts the projection of a given bbox.
sm.Convert(bbox []float64, to string) []float64

// Forward converts lon/lat values to 900913 x/y.
sm.Forward(ll []float64) []float64

// Inverse converts 900913 x/y values to lon/lat.
sm.Inverse(xy []float64) []float64

```
