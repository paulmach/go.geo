package clustering

import "math"

// A distanceSet is used to denormalize the minimum cluster distance in the set.
// Stores cluster-cluster distances, so MinDistance would be closest
// other cluster to the given cluster.
type distanceSet struct {
	MinDistance float64
	MinIndex    int
	Distances   map[int]float64
}

// newDistanceSet creates a new distance set.
func newDistanceSet() *distanceSet {
	return &distanceSet{
		MinDistance: math.MaxFloat64,
		Distances:   make(map[int]float64, 500), // adding 500 was a 20% performance win
	}
}

// Set sets the distance for a given index
// and updates the denormalized min distance if necessary.
// Returns true if the minimum has been updated.
func (ds *distanceSet) Set(index int, distance float64) bool {
	if d, ok := ds.Distances[index]; ok && index == ds.MinIndex && distance > d {
		// minimum index is being made greater, need to reset
		ds.Distances[index] = distance
		ds.reset()
		return true
	}

	ds.Distances[index] = distance
	if distance < ds.MinDistance {
		ds.MinDistance = distance
		ds.MinIndex = index
		return true
	}

	return false
}

// Delete removes a values for an index and resets denormalized minimum.
// Returns true if the minimum has been updated.
func (ds *distanceSet) Delete(index int) bool {
	delete(ds.Distances, index)
	if index == ds.MinIndex {
		ds.reset()
		return true
	}

	return false
}

// reset refinds the denormalized minimum. Can be used if manually updating values.
// not part of the official api.
func (ds *distanceSet) reset() {
	minDist := math.MaxFloat64
	minIndex := 0
	for i, d := range ds.Distances {
		if d < minDist {
			minDist = d
			minIndex = i
		}
	}

	ds.MinIndex = minIndex
	ds.MinDistance = minDist
}
