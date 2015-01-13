go.geo/clustering
=================

This package provides simple hierarchical clustering for points and paths with the 
[point\_clusting](point_clustering) and [path\_clustering](path_clustering)
packages, respectively.

Code overview:

* [**point_clusting**](point_clustering) Structs and methods for clustering points.
* [**path_clusting**](path_clusting) Structs and methods for clustering paths. 
	This is a still a work in progress.
* [**helpers**](helpers) Random functions for working with the results. 
	Stuff like prefiltering, point to cluster matching, etc.
* [**shared**](shared) Shared code for maintaining the next nearest pair data structure
	as well as caching of distances.

The code between the two libraries is very similar. It was duplicated to remove 
type assertions and allow for easier optimization for the specific types.
While ugly, it did result in a 5% performance gain, so my pain is your gain :)
