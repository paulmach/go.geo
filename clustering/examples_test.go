package clustering_test

import (
	"fmt"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/clustering"
)

func ExamplePointClustering() {
	pointers := []clustering.Pointer{
		&Event{Location: geo.NewPoint(1, 1)},
		&Event{Location: geo.NewPoint(2, 2)},
		&Event{Location: geo.NewPoint(5, 5)},
	}

	clusters := clustering.ClusterPointers(
		pointers,
		clustering.CentroidDistance{},
		2, // distance threshold, merge until clusters are at least this far apart
	)

	for i, c := range clusters {
		fmt.Printf("cluster %d:\n", i+1)
		for _, p := range c.Pointers {
			e := p.(*Event)
			fmt.Printf("   %+v\n", e)
		}
	}
	// Output:
	// cluster 1:
	//    &{Location:POINT(1 1)}
	//    &{Location:POINT(2 2)}
	// cluster 2:
	//    &{Location:POINT(5 5)}
}

func ExampleGeoPointClustering() {
	pointers := []clustering.Pointer{
		&Event{Location: geo.NewPoint(-122.548081, 37.905995)},
		&Event{Location: geo.NewPoint(-122.548091, 37.905987)},
		&Event{Location: geo.NewPoint(-122.54807, 37.905995)},
		&Event{Location: geo.NewPoint(-122.54807, 37.905995)},
		&Event{Location: geo.NewPoint(-122.54807, 37.905987)},
	}

	threshold := 1.0 // meter
	clusters := clustering.ClusterGeoPointers(
		pointers,
		threshold,
	)

	for i, c := range clusters {
		fmt.Printf("cluster %d:\n", i+1)
		for _, p := range c.Pointers {
			e := p.(*Event)
			fmt.Printf("   %+v\n", e)
		}
	}
	// Output:
	// cluster 1:
	//    &{Location:POINT(-122.548081 37.905995)}
	// cluster 2:
	//    &{Location:POINT(-122.548091 37.905987)}
	// cluster 3:
	//    &{Location:POINT(-122.54807 37.905995)}
	//    &{Location:POINT(-122.54807 37.905995)}
	//    &{Location:POINT(-122.54807 37.905987)}
}

// example of an object implementing the point_clusting.Pointer interface
type Event struct {
	Location *geo.Point
}

func (e *Event) CenterPoint() *geo.Point {
	return e.Location
}
