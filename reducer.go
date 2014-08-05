package geo

// A Reducer reduces a path using any simplification algorithm.
// It should return a copy of the path, not motify the original.
type Reducer interface {
	Reduce(*Path) *Path
}
