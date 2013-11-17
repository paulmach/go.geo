package geo

import (
	"fmt"
	"io"
	"math"
)

// Surface is kind of the 2d version of path.
type Surface struct {
	Bound         *Bound
	Width, Height int
	Grid          [][]float64 // x,y
}

func NewSurface(bound *Bound, width, height int) *Surface {
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
func (s *Surface) PointAt(x, y int) *Point {
	if x < 0 || x >= s.Width || y < 0 || y >= s.Height {
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
	if !s.Bound.Contains(point) {
		return &Point{}
	}

	var w1, w2 float64

	// find height and width
	w := (point[0] - s.Bound.sw[0]) / s.Bound.Width() * float64(s.Width-1)
	h := (point[1] - s.Bound.sw[1]) / s.Bound.Height() * float64(s.Height-1)

	xi := int(math.Floor(w))
	yi := int(math.Floor(h))

	// vertical
	h -= math.Floor(h)

	w1 = s.Grid[xi][yi]*w + s.Grid[xi+1][yi]*(1-w)
	w2 = s.Grid[xi][yi+1]*w + s.Grid[xi+1][yi+1]*(1-w)

	vertical := &Point{1, 0}
	vertical.Scale(w2 - w1)

	// horizontal
	w -= math.Floor(w)

	w1 = s.Grid[xi][yi]*w + s.Grid[xi][yi+1]*(1-w)
	w2 = s.Grid[xi][yi+1]*w + s.Grid[xi+1][yi+1]*(1-w)

	horizontal := &Point{1, 0}
	horizontal.Scale(w2 - w1)

	return &Point{}
}

// WriteOffFile writes an Object File Format representation of
// the surface to the writer provided. This is for viewing
// in MeshLab or something like that. You should close the
// writer yourself after this function returns.
// http://segeval.cs.princeton.edu/public/off_format.html
func (s *Surface) WriteOffFile(w io.Writer) {
	facesCount := 0
	faces := ""

	for i := 0; i < s.Width-1; i += 1 {
		for j := i % 2; j < s.Height-1; j += 2 {
			faces += fmt.Sprintf("4 %d %d %d %d\n", i*s.Height+j, i*s.Height+j+1, (i+1)*s.Height+j+1, (i+1)*s.Height+j)
			facesCount++
		}
	}

	w.Write([]byte("OFF\n"))
	w.Write([]byte(fmt.Sprintf("%d %d 0\n", s.Height*s.Width, facesCount)))

	// vertexes
	for i := 0; i < s.Width; i++ {
		for j := 0; j < s.Height; j++ {
			p := s.PointAt(i, j)
			w.Write([]byte(fmt.Sprintf("%f %f %f\n", p[0], p[1], s.Grid[i][j])))
		}
	}

	w.Write([]byte(faces))
}
