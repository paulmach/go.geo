package geo

import (
	"fmt"
	"math"
)

// A Bound represents and enclosed "box" in the 2d Euclidean or Cartesian plane.
// It does not know anything about the anti-meridian.
type Bound struct {
	sw, ne *Point
}

func NewBound(east, west, north, south float64) *Bound {
	return &Bound{
		sw: &Point{math.Min(east, west), math.Min(north, south)},
		ne: &Point{math.Max(east, west), math.Max(north, south)},
	}
}

func NewBoundFromPoints(corner, oppositeCorner *Point) *Bound {
	b := &Bound{
		sw: corner.Clone(),
		ne: corner.Clone(),
	}

	b.Extend(oppositeCorner)
	return b
}

// Extend grows the bound to include the new point
func (b *Bound) Extend(point *Point) *Bound {

	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	b.sw.SetX(math.Min(b.sw.X(), point.X()))
	b.ne.SetX(math.Max(b.ne.X(), point.X()))

	b.sw.SetY(math.Min(b.sw.Y(), point.Y()))
	b.ne.SetY(math.Max(b.ne.Y(), point.Y()))

	return b
}

// Contains figures out if the point is within the bounds.
// On the boundary is considered in.
func (b *Bound) Contains(point *Point) bool {

	if point.Y() < b.sw.Y() || b.ne.Y() < point.Y() {
		return false
	}

	if point.X() < b.sw.X() || b.ne.X() < point.X() {
		return false
	}

	return true
}

// Intersects determines if two bounds intersect.
// True if they are touching.
func (b *Bound) Intersects(bound *Bound) bool {
	if bound.Contains(b.sw) || bound.Contains(b.ne) ||
		bound.Contains(b.SouthEast()) || bound.Contains(b.NorthWest()) {
		return true
	}

	// now check the completely inside case, only one consition required
	if b.Contains(bound.sw) {
		return true
	}

	return false
}

func (b *Bound) Center() *Point {
	p := &Point{}
	p.SetX((b.ne.X() + b.sw.X()) / 2.0)
	p.SetY((b.ne.Y() + b.sw.Y()) / 2.0)

	return p
}

// Expands in all directions by the amount given. The amount must be
// in the units of the bounds. Techinally one can pad with negative value,
// but no error checking is done.
func (b *Bound) Pad(amount float64) *Bound {
	b.sw.SetX(b.sw.X() - amount)
	b.sw.SetY(b.sw.Y() - amount)

	b.ne.SetX(b.ne.X() + amount)
	b.ne.SetY(b.ne.Y() + amount)

	return b
}

// Height returns just the difference in the points' Y/Latitude
func (b *Bound) Height() float64 {
	return b.ne.Y() - b.sw.Y()
}

// Width returns just the difference in the points' X/Longitude
func (b *Bound) Width() float64 {
	return b.ne.X() - b.sw.X()
}

// GeoHeight returns the approximate height in meters.
// Only applies if the data is Lat/Lng degrees.
func (b *Bound) GeoHeight() float64 {
	return 111131.75 * b.Height()
}

// GeoWidth returns the approximate width in meters.
// Only applies if the data is Lat/Lng degrees.
func (b *Bound) GeoWidth(haversine ...bool) float64 {
	c := b.Center()

	A := &Point{b.sw[0], c[1]}
	B := &Point{b.ne[0], c[1]}

	return A.GeoDistanceFrom(B, yesHaversine(haversine))
}

func (b *Bound) SouthWest() *Point { return b.sw.Clone() }
func (b *Bound) NorthEast() *Point { return b.ne.Clone() }

func (b *Bound) SouthEast() *Point {
	newP := &Point{}
	newP.SetLat(b.sw.Lat()).SetLat(b.ne.Lng())
	return newP
}

func (b *Bound) NorthWest() *Point {
	newP := &Point{}
	newP.SetLat(b.ne.Lat()).SetLat(b.sw.Lng())
	return newP
}

func (b *Bound) Empty() bool {
	return b.sw.Equals(b.ne)
}

func (b *Bound) Equals(c *Bound) bool {
	if b.sw.Equals(c.sw) && b.ne.Equals(c.ne) {
		return true
	}

	return false
}

func (b *Bound) Clone() *Bound {
	return NewBoundFromPoints(b.sw, b.ne)
}

func (b *Bound) String() string {
	return fmt.Sprintf("[[%f, %f], [%f, %f]]", b.sw.X(), b.ne.X(), b.sw.Y(), b.ne.Y())
}
