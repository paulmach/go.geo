go.geo/reducers
===============

This package implements several reducers that simplify a geo.Path object.
See the [reducers godoc](http://godoc.org/github.com/paulmach/go.geo/reducers) for more information.

Note: all these methods **create a new path** and do not modify the input path.

Currently implemented:

* [Douglas-Peucker](#dp)
* [Visvalingam](#vis)
* [Radial](http://psimpl.sourceforge.net/radial-distance.html)

Performance
-----------

These reducers are optimized and performance is comparible to libraries in other languages.

go get github.com/paulmach/go.geo/reducers
go test github.com/paulmach/go.geo/reducers -bench .


<a name="dp"></a>Douglas-Peucker
--------------------------------

Probably the most popular simplification algorithm around. For algorithm details, see
[wikipedia](http://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm).

Usage: 

	originalPath := geo.NewPath()
	reducedPath := reducers.DouglasPeucker(originalPath, threshold)

	// the index map method can be used to figure out which of the points were kept
	reducedPath, indexMap := reducers.DouglasPeuckerIndexMap(originalPath, threshold)

	for i, v := range indexMap {
		reducedPath.GetAt(i) == originalPath.GetAt(v)
	}

	// to chain reducers and combine their index maps use MergeIndexMaps
	p1, im1 := reducers.RadialIndexMap(path, meters) 
	reducedPath, im2 := reducers.DouglasPeuckerIndexMap(p1, threshold)
	indexMap := MergeIndexMaps(im1, im2)

<a name="vis"></a>Visvalingam
-----------------------------

See Mike Bostock's explanation for 
[algorithm details](http://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm).

Usage: 

	originalPath := geo.NewPath()
	reducedPath := reducers.DouglasPeucker(originalPath, threshold)

	// will remove all whose triangle is smaller than `threshold`
	reducedPath := reducers.VisvalingamThreshold(path, threshold)

	// will remove points until there are only `toKeep` points left.
	reducedPath := reducers.VisvalingamKeep(path, toKeep)

	// One can also combine the parameters.
	// This will continue to remove points until 
	//  - there are no more below the threshold,
	//  - or the new path is of length `toKeep`
	reducedpath := reducers.Visvalingam(path, threshold, toKeep)

<a name="radial"></a>Radial
---------------------------

Radial reduces the path by removing points that are close together.
A full [algorithm description](http://psimpl.sourceforge.net/radial-distance.html).

Usage: 

	originalPath := geo.NewPath()

	// this method uses a Euclidean distance measure.
	reducedPath := reducers.Radial(path, meters)

	// if the points are in the lng/lat space Radial Geo will 
	// compute the geo distance between the coordinates.
	reducedPath := reducers.RadialGeo(path, meters)
