package clustering

import "math"

// State represents the state of the hierarchical clustering and manages
// the updates of the distance sets.
type state struct {
	Distances    []*distanceSet
	DistanceFunc func(a, b int) float64
}

// ResetDistances makes sure the distance map is up to date given the recent merge of clusters.
func (s *state) ResetDistances(into, from int) {
	// since the center of into changed, need to update the distance to anything linked to this one.
	for k := range s.Distances[into].Distances {
		if k == into {
			continue
		}

		dist := s.DistanceFunc(into, k)

		s.Distances[into].Set(k, dist)
		s.Distances[k].Set(into, dist)
	}

	// we are merging from into into.
	// so any links with from need to be links with into.
	// any links into from need to be deleted
	// we no longer need the distance info of from
	for k := range s.Distances[from].Distances {
		if k == from || k == into {
			continue
		}

		dist := s.DistanceFunc(into, k)

		s.Distances[into].Set(k, dist)
		s.Distances[k].Set(into, dist)

		s.Distances[k].Delete(from)
	}

	s.Distances[from] = nil
	s.Distances[into].Delete(from)
}

// MinDistance returns the link with minimum distance.
// a is the index stored on the DistanceSet, b is the index of the smallest values.
func (s *state) MinDistance() (a, b int, dist float64) {
	dist = math.MaxFloat64

	for i, ds := range s.Distances {
		if ds == nil {
			continue
		}

		if ds.MinDistance < dist {
			dist = ds.MinDistance

			a = i
			b = ds.MinIndex
		}
	}

	return
}
