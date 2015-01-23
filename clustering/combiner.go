package clustering

import "math"

// A Combiner is something that can be combined.
type Combiner interface {
	Combine(c Combiner) Combiner
	DistanceFromCombiner(c Combiner) float64
}

// ClusterCombiners will do a simple hierarchical of the combiners.
// It will modify the input slice as things will be combined into each other.
func ClusterCombiners(combiners []Combiner, threshold float64) []Combiner {
	if len(combiners) < 2 {
		return combiners
	}

	s := &state{
		Distances: initializeCombinerDistances(combiners, threshold),
		DistanceFunc: func(a, b int) float64 {
			return combiners[a].DistanceFromCombiner(combiners[b])
		},
	}

	// successively merge
	for i := 1; i < len(combiners); i++ {
		lower, higher, dist := s.MinDistance()
		if dist > threshold {
			break
		}

		// merge these two
		combiners[lower] = combiners[lower].Combine(combiners[higher])
		s.ResetDistances(lower, higher)
		combiners[higher] = nil
	}

	last := len(combiners) - 1
	for i := range combiners {
		if combiners[i] == nil {
			combiners[i] = combiners[last]
			last--
		}
	}

	return combiners[:last+1]
}

func initializeCombinerDistances(
	combiners []Combiner,
	threshold float64,
) []*distanceSet {

	// initialize distances
	distances := make([]*distanceSet, len(combiners))

	for i := range combiners {

		if distances[i] == nil {
			distances[i] = newDistanceSet()
		}
		distances[i].Set(i, math.MaxInt32)

		for j := i + 1; j < len(combiners); j++ {
			// TODO: better filtering here we don't have a literal n^2 situation.
			dist := combiners[i].DistanceFromCombiner(combiners[j])
			if dist < 5*threshold {
				distances[i].Set(j, dist)

				if distances[j] == nil {
					distances[j] = newDistanceSet()
				}
				distances[j].Set(i, dist)
			} else {
				// greater than a big threshold, so pass
			}
		}
	}

	return distances
}
