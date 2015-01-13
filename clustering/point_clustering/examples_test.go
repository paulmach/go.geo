package point_clustering_test

import (
	"fmt"

	"github.com/paulmach/go.geo"
	"github.com/paulmach/go.geo/clustering/point_clustering"
)

func ExamplePointClustering() {
	pointers := []point_clustering.Pointer{
		&Event{Location: geo.NewPoint(1, 1)},
		&Event{Location: geo.NewPoint(2, 2)},
		&Event{Location: geo.NewPoint(5, 5)},
	}

	clusters := point_clustering.New(
		2, // distance threshold, merge until clusters are at least this far apart
		point_clustering.CentroidDistance{},
	).Cluster(pointers)

	for i, c := range clusters {
		fmt.Printf("cluster %d:\n", i+1)
		for _, p := range c.Pointers {
			e := p.(*Event)
			fmt.Printf("   %+v\n", e)
		}
	}
	// Output:
	// cluster 1:
	//    &{Location:[1.000000, 1.000000]}
	//    &{Location:[2.000000, 2.000000]}
	// cluster 2:
	//    &{Location:[5.000000, 5.000000]}
}

func ExampleGeoPointClustering() {
	pointers := []point_clustering.Pointer{
		&Event{Location: geo.NewPoint(-122.548081, 37.905995)},
		&Event{Location: geo.NewPoint(-122.548091, 37.905987)},
		&Event{Location: geo.NewPoint(-122.54807, 37.905995)},
		&Event{Location: geo.NewPoint(-122.54807, 37.905995)},
		&Event{Location: geo.NewPoint(-122.54807, 37.905987)},
	}

	threshold := 1.0 // meter
	clusters := point_clustering.NewGeoProjectedClustering(threshold).Cluster(pointers)

	for i, c := range clusters {
		fmt.Printf("cluster %d:\n", i+1)
		for _, p := range c.Pointers {
			e := p.(*Event)
			fmt.Printf("   %+v\n", e)
		}
	}
	// Output:
	// cluster 1:
	//    &{Location:[-122.548081, 37.905995]}
	// cluster 2:
	//    &{Location:[-122.548091, 37.905987]}
	// cluster 3:
	//    &{Location:[-122.548070, 37.905995]}
	//    &{Location:[-122.548070, 37.905995]}
	//    &{Location:[-122.548070, 37.905987]}
}

// example of an object implementing the point_clusting.Pointer interface
type Event struct {
	Location *geo.Point
}

func (e *Event) CenterPoint() *geo.Point {
	return e.Location
}
