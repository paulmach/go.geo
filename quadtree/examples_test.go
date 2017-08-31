package quadtree_test

import (
	"fmt"
	"math/rand"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/quadtree"
)

func ExampleQuadtreeFind() {
	r := rand.New(rand.NewSource(42)) // to make things reproducible

	qt := quadtree.New(geo.NewBound(0, 1, 0, 1))

	// insert 1000 random points
	for i := 0; i < 1000; i++ {
		qt.Insert(geo.NewPoint(r.Float64(), r.Float64()))
	}

	nearest := qt.Find(geo.NewPoint(0.5, 0.5))
	fmt.Printf("nearest: %+v\n", nearest)

	// Output:
	// nearest: POINT(0.4930591659434973 0.5196585530161364)
}

func ExampleQuadtreeFindMatching() {
	r := rand.New(rand.NewSource(42)) // to make things reproducible

	type dataPoint struct {
		geo.Pointer
		visible bool
	}

	qt := quadtree.New(geo.NewBound(0, 1, 0, 1))

	// insert 100 random points
	for i := 0; i < 100; i++ {
		qt.Insert(dataPoint{geo.NewPoint(r.Float64(), r.Float64()), false})
	}

	qt.Insert(dataPoint{geo.NewPoint(0, 0), true})

	nearest := qt.FindMatching(
		geo.NewPoint(0.5, 0.5),
		func(p geo.Pointer) bool { return p.(dataPoint).visible },
	)

	fmt.Printf("nearest: %+v\n", nearest)

	// Output:
	// nearest: {Pointer:POINT(0 0) visible:true}
}

func ExampleQuadtreeInBound() {
	r := rand.New(rand.NewSource(52)) // to make things reproducible

	qt := quadtree.New(geo.NewBound(0, 1, 0, 1))

	// insert 1000 random points
	for i := 0; i < 1000; i++ {
		qt.Insert(geo.NewPoint(r.Float64(), r.Float64()))
	}

	bounded := qt.InBound(geo.NewBound(0.5, 0.5, 0.5, 0.5).Pad(0.05))
	fmt.Printf("in bound: %v\n", len(bounded))
	// Output:
	// in bound: 10
}
