// Package primesingo provides efficient algorithms for listing prime numbers.
//
// The main entry point is [Sieve], which returns all primes up to a given
// bound using a parallel segmented Sieve of Eratosthenes. The implementation
// is designed for educational clarity while remaining performant: it uses
// bit-packed storage, cache-sized segments, and CPU-bounded parallelism.
//
// Basic usage:
//
//	ps := primesingo.Sieve(100)
//	fmt.Println(ps) // [2 3 5 7 11 13 ... 97]
//
// For very large N, prefer streaming approaches over materializing the full
// list — the result slice itself becomes the dominant memory cost.

package primesingo
