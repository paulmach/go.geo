go.geo
======

Go.geo is a geometry/geography library in [Go](http://golang.org). Its purpose is to allow for
basic point, line and path operations on the server side or while scripting. The primary use
case being GIS geometry manipulation on the server side vs. in the browser using javascript.
This may be motivated by memory, computation time or data privacy constraints.
All objects are defined in a 2D context.

#### To install
	
	go get github.com/paulmach/go.geo

#### To use, imports as package name `geo`:

	import "github.com/paulmach/go.geo"

<br />
[![Build Status](https://travis-ci.org/paulmach/go.geo.png?branch=master)](https://travis-ci.org/paulmach/go.geo)
&nbsp; &nbsp;
[![Coverage Status](https://coveralls.io/repos/paulmach/go.geo/badge.png?branch=master)](https://coveralls.io/r/paulmach/go.geo?branch=master)
&nbsp; &nbsp;
[![Godoc Reference](https://godoc.org/github.com/paulmach/go.geo?status.png)](https://godoc.org/github.com/paulmach/go.geo)

## Exposed objects

* **Point** represents a 2D location, x/y or lng/lat.
	It also supports some vector functions like add, scale, etc.
	It's up to the programmer to know if the data is a lng/lat location, 
	projection of that point, or a vector.
* **Line** represents the shortest distance between two points in Euclidean space. 
	In many cases the path object is more useful.
* **Path** represents a set of points making up a path in the 2D plane.
	Functions for converting to/from
	[Google's polyline encoding](https://developers.google.com/maps/documentation/utilities/polylinealgorithm) are included.
* **Bound** represents a rectangular 2D area defined by North, South, East, West values.
	Computable for Line and Path objects, used by the Surface object.
* **Surface** is used to assign values to points in a 2D area, such as elevation.

### Reducers

The reducers sub-package includes implementations for 
[Douglas-Peucker](http://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm),
[Visvalingam](http://bost.ocks.org/mike/simplify/) 
and
[Radial](http://psimpl.sourceforge.net/radial-distance.html) 
polyline reduction algorithms. See the [reducers godoc](http://godoc.org/github.com/paulmach/go.geo/reducers) for more information.

## Library conventions

There are two big conventions that developers should be aware of:
**functions are chainable** and **operations modify the original object.**
For example:

	p := geo.NewPoint(0, 0)
	p.SetX(10).Add(geo.NewPoint(10, 10))
	p.Equals(geo.NewPoint(20, 10))  // == true

If you want to create a copy, all objects support the `Clone` method.

	p1 := geo.NewPoint(10, 10)
	p2 := p1.SetY(20)
	p1.Equals(p2) // == true, in this case p1 and p2 point to the same memory

	p2 := p1.Clone().SetY(30)
	p1.Equals(p2) // == false

These conventions put extra load on the programmer,
but tests showed that making a copy every time was significantly slower.
So, **remember to explicitly Clone() your objects**.

## Performance

As this code is meant to act as a core library, performance is very important. 
There is a good set of benchmarks covering the core objects and efforts have been
made to optimize them. To run:

	go get github.com/paulmach/go.geo
	go test github.com/paulmach/go.geo -bench .

To benchmark the reducers run:

	go get github.com/paulmach/go.geo/reducers
	go test github.com/paulmach/go.geo/reducers -bench .

## Examples

The [GoDoc Documentation](https://godoc.org/github.com/paulmach/go.geo) provides a very readable list
of exported functions. Below are a few usage examples.

### Encode/Decode polyline path

	// lng/lat data, in this case, is encoded at 6 decimal place precision
	path := geo.Decode("smsqgAtkxvhFwf@{zCeZeYdh@{t@}BiAmu@sSqg@cjE", 1e6)
	
	// reduce using the Douglas Peucker line reducer from the reducers sub-package.
	// Note the threshold distance is in the coordinates of the points,
	// which in this case is degrees.
	reducedPath := reducers.DouglasPeucker(path, 1.0e-5)
	
	// encode with the default/typical 5 decimal place precision
	encodedString := reducedPath.Encode() 

### Path, line intersection

	path := geo.NewPath()
	path.Push(geo.NewPoint(0, 0))
	path.Push(geo.NewPoint(1, 1))

	line := geo.NewLine(geo.NewPoint(0, 1), geo.NewPoint(1, 0))

	// intersects does a simpler check for yes/no
	if path.Intersects(line) {
		// intersection will return the actual points and places on intersection
		points, segments := path.Intersection(line)

		for i, _ := range points {
			log.Printf("Intersection %d at %v with path segment %d", i, points[i], segments[i][0])
		}
	}

## Surface

A surface object is defined by a bound (lng/lat georegion for example) and a width and height 
defining the number of discrete points in the bound. This allows for access such as:

	surface.Grid[x][y]         // the value at a location in the grid
	surface.GetPoint(x, y)     // the point, which will be in the space as surface.bound,
	                           // corresponding to surface.Grid[x][y]
	surface.ValueAt(*Point)    // the bi-linearly interpolated grid value for any point in the bounds
	surface.GradientAt(*Point) // the gradient of the surface a any point in the bounds,
	                           // returns a point object which should be treated as a vector

A couple things about how the bound area is discretized in the grid:
 
	* surface.Grid[0][0]
		corresponds to the surface.Bound.SouthWest() location, or bottom left corner or the bound
	* surface.Grid[0][surface.Height-1]
		corresponds to the surface.Bound.NorthWest() location,
		the extreme points in the grid are on the edges of the bound

While these conventions are useful, they are different.
If you're using this object, your feedback on these choices would be appreciated.

## Contributing

While this project started as the core of [Slide](https://github.com/paulmach/slide) it's now being used in many place.
So, if you have features you'd like to add or improvements to make, please submit a pull request.
