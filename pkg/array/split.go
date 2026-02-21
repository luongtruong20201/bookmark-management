// Package array provides generic utilities for slice operations, such as
// splitting slices into fixed-size chunks.
package array

// SplitIntoChunks splits the given slice into sub-slices of at most size elements each.
// The last chunk may contain fewer than size elements when len(objects) is not
// divisible by size. Returns a single empty chunk for an empty input slice.
// size must be positive; behavior is undefined for size <= 0.
func SplitIntoChunks[T any](objects []T, size int) [][]T {
	chunks := make([][]T, 0, (len(objects)+size-1)/size)

	for size < len(objects) {
		objects, chunks = objects[size:], append(chunks, objects[0:size:size])
	}
	chunks = append(chunks, objects)
	return chunks
}
