// Package collection provides fluent, generic wrappers for working with slices and maps.
//
// The core types are:
//
//   - [Collection] — a generic wrapper around a slice with a rich set of chainable methods
//     for filtering, sorting, transforming, and aggregating data.
//   - [MapCollection] — an ordered map with a fluent API for key-value operations,
//     preserving insertion order.
//   - [LazyCollection] — a lazily-evaluated sequence backed by [iter.Seq], allowing
//     efficient pipeline processing of large or infinite datasets.
//
// All three types support conversion to Go 1.23+ iterators via their Iter methods,
// enabling use with range-over-func loops.
//
// # Sub-packages
//
// Standalone utility functions are available in sub-packages:
//
//   - [github.com/gocanto/collection/arr] — generic slice helpers (filter, sort, partition, etc.)
//   - [github.com/gocanto/collection/kv]  — map helpers with dot-notation support for nested maps
package collection
