package primesingo

import "math"

// Sieve returns all prime numbers less than or equal to n.
//
// The returned slice is sorted in ascending order. For n < 2, Sieve
// returns nil, since there are no primes below 2.
//
// Sieve is safe for concurrent use; it allocates its own working memory
// per call and shares no mutable state.
func Sieve(n uint64) []uint64 {
	if n < 2 {
		return nil
	}

	composite := make([]bool, n+1)

	// Marking phase.
	for p := uint64(2); p*p <= n; p++ {
		if composite[p] {
			continue
		}
		for m := p * p; m <= n; m += p {
			composite[m] = true
		}
	}

	// Collection phase: pre-allocate using π(N) ≈ N/ln(N) from the
	// Prime Number Theorem. The 1.1 multiplier covers the small-N
	// underestimate of the approximation; the +16 floor handles very
	// small N where N/ln(N) is numerically meaningless.
	estimate := 16
	if n >= 10 {
		estimate = int(1.1*float64(n)/math.Log(float64(n))) + 16
	}
	primes := make([]uint64, 0, estimate)

	for i := uint64(2); i <= n; i++ {
		if !composite[i] {
			primes = append(primes, i)
		}
	}

	return primes
}
