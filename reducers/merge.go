package reducers

// MergeIndexMaps merges two index maps for use when chaining reducers.
// For example, to radially reduce and then DP, merge the index maps with this function
// to get a map from the original to the final path.
func MergeIndexMaps(map1, map2 []int) []int {
	result := make([]int, len(map2))
	for i, v := range map2 {
		result[i] = map1[v]
	}

	return result
}
