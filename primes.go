package primesingo

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

	// The sieve itself: composite[i] tells us whether we've determined
	// i is composite. The zero value of bool is false, so a freshly
	// allocated slice already means "everything is presumed prime" —
	// we get the initialization for free without an O(n) setup pass.
	//
	// Why name it `composite` instead of `isPrime` or `sieve`?
	// The marking loop only ever *sets* entries to true (never false).
	// Naming the slice for what `true` means makes the code read
	// naturally: "if composite[m] is already true, skip; otherwise
	// mark composite[m] = true."
	composite := make([]bool, n+1)

	// Outer loop: try each candidate prime p from 2 up to √n.
	// We use p*p <= n instead of p <= sqrt(n) to avoid floating-point
	// rounding bugs — same reason we used d*d <= n in trial division.
	//
	// Why stop at √n? If a number c ≤ n is composite, it factors as
	// c = a·b with a ≤ b. Then a² ≤ a·b = c ≤ n, so a ≤ √n. The smallest
	// prime factor of any composite ≤ n is ≤ √n, so once we've processed
	// all primes up to √n, every composite has been marked.
	for p := uint64(2); p*p <= n; p++ {
		if composite[p] {
			continue // p was marked by a smaller prime; not a prime itself
		}

		// p is prime. Mark all its multiples in [p², n] as composite.
		//
		// Why start at p² and not 2p? Every multiple of p smaller than
		// p² has the form p·k for some k < p. That k must have a prime
		// factor q ≤ k < p, and when we processed q earlier, we already
		// marked p·k (which is also a multiple of q). Starting at p²
		// skips the redundant work.
		for m := p * p; m <= n; m += p {
			composite[m] = true
		}
	}

	// Collection phase: walk the sieve, collect indices that survived.
	// We don't know exactly how many primes we'll find, so we let
	// append grow the slice. Pre-allocation comes in a later step.
	var primes []uint64
	for i := uint64(2); i <= n; i++ {
		if !composite[i] {
			primes = append(primes, i)
		}
	}

	return primes
}
