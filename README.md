go.geo
======

Go.geo is a geometry/geography library in [Go](http://golang.org). Its purpose is to allow for
basic point, line and path operations in server side operations or scripting. The primary use
case being GIS geometry manipulation on the server side vs. in the browser using javascript.
This may be motivated by memory, computation time or data privacy constraints.
All objects are defined in a 2D context.

Development is currently focused on my needs. But pull requests are welcome.

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
	It's up to the programmer to know if the data is lng/lat location or a 
	projection of that point, or whatever.
* **Line** represents the shortest distance between two points in Euclidean space. 
	In many cases the path object is more useful.
* **Path** represents a set of points representing a path in the 2D plain. 
	Function for converting to/from
	[Google's polyline encoding](https://developers.google.com/maps/documentation/utilities/polylinealgorithm) are included.
* **Bound** represents a rectangular 2d area defined by North, South, East, West values.
* **Surface** is used to assign values to points in a 2D area, such as elevation.
	*This object is still being developed and is experimental*

## Library conventions

There are two big conventions that developers should be aware off; 
**functions are chainable** and **operations modify the original object.**
For example:

	p := geo.NewPoint(0, 0)
	p.SetX(10).Add(geo.NewPoint(10, 10))
	p.Equals(geo.NewPoint(20, 10))  // == true

If you want to create a copy, all the objects support he `Clone` method.

	p1 := geo.NewPoint(10, 10)
	p2 := p1.SetY(20)
	p1.Equals(p2) // == true, in this case p1 and p2 point to the same memory

	p2 := p1.Clone().SetY(30)
	p1.Equals(p2) // == false

I'm not sure if modifying the original object is the right choice for this library,
but some simple tests showed that making a copy each time was significantly slower.
So, **remember to Clone() your objects**.

## Examples

// TODO
