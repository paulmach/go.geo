package geo

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

// Surface is kind of the 2d version of path.
type Surface struct {
	Bound         *Bound
	Width, Height uint32

	// represents the underlying data, as [x][y]
	// where x in [0:Width] and y in [0:Height]
	Grid [][]float64 // x,y
}

func NewSurface(bound *Bound, width, height uint32) *Surface {
	s := &Surface{
		Bound:  bound.Clone(),
		Width:  width,
		Height: height,
	}

	s.Grid = make([][]float64, width)
	points := make([]float64, width*height)

	for i := range s.Grid {
		s.Grid[i], points = points[:height], points[height:]
	}

	return s
}

// PointAt returns the point corresponding to this grid coordinate
// given the size and bounds of the surface.
// x in [0, s.Width()-1]
// y in [0, s.Height()-1]
func (s *Surface) PointAt(x, y uint32) *Point {
	if x >= s.Width || y >= s.Height {
		return nil
	}

	w := float64(x) / float64(s.Width-1)
	h := float64(y) / float64(s.Height-1)

	p := s.Bound.sw.Clone()

	p[0] += w * s.Bound.Width()
	p[1] += h * s.Bound.Height()

	return p
}

// ValueAt returns the bi-linearly interpolated value for
// the given point. Returns 0 if the point is out of surface bounds
// TODO: cleanup and optimize this code
func (s *Surface) ValueAt(point *Point) float64 {
	var w1, w2 float64

	if !s.Bound.Contains(point) {
		return 0
	}

	// find height and width
	w := (point[0] - s.Bound.sw[0]) / s.Bound.Width() * float64(s.Width-1)
	h := (point[1] - s.Bound.sw[1]) / s.Bound.Height() * float64(s.Height-1)

	xi := int(math.Floor(w))
	yi := int(math.Floor(h))

	xi1 := int(math.Ceil(w))
	yi1 := int(math.Ceil(h))

	w -= math.Floor(w)
	h -= math.Floor(h)

	w1 = s.Grid[xi][yi]*(1-w) + s.Grid[xi1][yi]*w
	w2 = s.Grid[xi][yi1]*(1-w) + s.Grid[xi1][yi1]*w

	return w1*(1-h) + w2*h
}

func (s *Surface) GradientAt(point *Point) *Point {
	delta := s.Bound.Width() / float64(s.Width-1) / 5.0

	if !s.Bound.Clone().Pad(delta).Contains(point) {
		return &Point{}
	}

	// horizontal
	x1 := s.ValueAt(point.Clone().Add(NewPoint(-delta, 0)))
	x2 := s.ValueAt(point.Clone().Add(NewPoint(delta, 0)))

	horizontal := NewPoint((x1-x2)/(2*delta), 0)

	// vertical
	delta = s.Bound.Height() / float64(s.Height-1)
	y1 := s.ValueAt(point.Clone().Add(NewPoint(0, -delta)))
	y2 := s.ValueAt(point.Clone().Add(NewPoint(0, delta)))

	return horizontal.SetY((y1 - y2) / (2 * delta))
}

// WriteOffFile writes an Object File Format representation of
// the surface to the writer provided. This is for viewing
// in MeshLab or something like that. You should close the
// writer yourself after this function returns.
// http://segeval.cs.princeton.edu/public/off_format.html
func (s *Surface) WriteOffFile(w io.Writer) {
	var i, j uint32

	facesCount := 0
	var faces bytes.Buffer

	for i = 0; i < s.Width-1; i += 1 {
		for j := i % 2; j < s.Height-1; j += 2 {
			face := fmt.Sprintf("4 %d %d %d %d\n", i*s.Height+j, i*s.Height+j+1, (i+1)*s.Height+j+1, (i+1)*s.Height+j)
			faces.WriteString(face)
			facesCount++
		}
	}

	w.Write([]byte("OFF\n"))
	w.Write([]byte(fmt.Sprintf("%d %d 0\n", s.Height*s.Width, facesCount)))

	// vertexes
	for i = 0; i < s.Width; i++ {
		for j = 0; j < s.Height; j++ {
			p := s.PointAt(i, j)
			w.Write([]byte(fmt.Sprintf("%.8f %.8f %.8f\n", p[0], p[1], s.Grid[i][j])))
		}
	}

	w.Write(faces.Bytes())
}
