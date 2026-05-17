package primesingo

import (
	"reflect"
	"testing"
)

// TestSieve_SmallCases verifies behavior on the boundary values where
// the algorithm has to handle "no primes exist" gracefully. These cases
// often have bugs because developers focus on the interesting path.
func TestSieve_SmallCases(t *testing.T) {
	tests := []struct {
		name string
		n    uint64
		want []uint64
	}{
		{"zero returns nil", 0, nil},
		{"one returns nil", 1, nil},
		{"two returns [2]", 2, []uint64{2}},
		{"ten returns first four primes", 10, []uint64{2, 3, 5, 7}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Sieve(tc.n)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Sieve(%d) = %v, want %v", tc.n, got, tc.want)
			}
		})
	}
}

// TestSieve_AllReturnedNumbersArePrime verifies that every number Sieve
// returns is actually prime, using trial division as an independent check.
// This catches the "false positive" failure mode: a sieve bug that leaves
// composite numbers unmarked.
//
// We use a small bound because trial division is O(n√n) total, which
// gets slow fast. The point is correctness, not speed.
//
//goland:noinspection GrazieInspectionRunner
func TestSieve_AllReturnedNumbersArePrime(t *testing.T) {
	const bound = 1000

	got := Sieve(bound)

	// Sanity check: π(1000) = 168. If we got far fewer, the rest of this
	// test is meaningless — fail loudly instead of passing vacuously.
	// The threshold is deliberately loose (100, not 168) so this guard
	// doesn't duplicate the precise count check in TestSieve_NoPrimesMissing.
	if len(got) < 100 {
		t.Fatalf("Sieve(%d) returned only %d primes — too few to validate",
			bound, len(got))
	}

	for _, p := range got {
		if !isPrimeTrialDivision(p) {
			t.Errorf("Sieve returned %d, which is not prime", p)
		}
	}
}

// TestSieve_NoPrimesMissing verifies that Sieve doesn't accidentally
// skip primes — the "false negative" failure mode. We compare counts
// against known values of π(n), the prime-counting function.
//
// π(n) values from OEIS sequence A000720.
func TestSieve_NoPrimesMissing(t *testing.T) {
	tests := []struct {
		n         uint64
		wantCount int
	}{
		{10, 4},     // 2, 3, 5, 7
		{100, 25},   // π(100) = 25
		{1000, 168}, // π(1000) = 168
		{10000, 1229},
	}

	for _, tc := range tests {
		got := Sieve(tc.n)
		if len(got) != tc.wantCount {
			t.Errorf("Sieve(%d) returned %d primes, want %d",
				tc.n, len(got), tc.wantCount)
		}
	}
}

// TestSieve_LargeN exercises the algorithm at a scale where bugs that
// only appear with large inputs would surface (e.g., overflow, capacity
// estimation errors, segment boundary issues once we add segments).
//
// Skipped with -short, since this takes seconds rather than milliseconds.
// Run locally with `go test -short` for fast feedback; CI runs the full
// suite without -short before merging.
func TestSieve_LargeN(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large-N test in -short mode")
	}

	const bound = 10_000_000
	const wantCount = 664579 // π(10⁷) from OEIS A006880

	got := Sieve(bound)
	if len(got) != wantCount {
		t.Errorf("Sieve(%d) returned %d primes, want %d",
			bound, len(got), wantCount)
	}
}

// isPrimeTrialDivision is the reference implementation: dead simple,
// obviously correct, slow. We use it only in tests to validate the
// real algorithm. Production code would never call this.
func isPrimeTrialDivision(n uint64) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	// Only check odd divisors up to √n. Same √n logic as the sieve itself —
	// any composite n has a factor ≤ √n. Using d*d <= n instead of
	// d <= sqrt(n) avoids floating-point rounding bugs near integer
	// boundaries: math.Sqrt returns a float64 approximation, and casting
	// back can be off by one, causing missed divisors.
	for d := uint64(3); d*d <= n; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}
