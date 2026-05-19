package primesingo

import "testing"

// BenchmarkSieve measures the wall-clock cost of Sieve at several
// scales. Go's benchmark runner repeats the inner loop b.N times and
// reports nanoseconds per operation; b.N is chosen automatically by
// the framework to make the total runtime statistically meaningful
// (~1 second by default).
//
// We use sub-benchmarks (b.Run) for the same reason we use sub-tests:
// each scale gets its own name in the output, and you can target a
// single scale with -bench='BenchmarkSieve/N=10000$'.
func BenchmarkSieve(b *testing.B) {
	scales := []uint64{
		1_000,
		10_000,
		100_000,
		1_000_000,
	}

	for _, n := range scales {
		// %d in the name is just for readability in benchmark output.
		// The "N=" prefix makes the column easier to scan.
		b.Run(fmtName(n), func(b *testing.B) {
			// b.ResetTimer is unnecessary here (we have no setup
			// before the loop), but it's a habit worth building for
			// benchmarks where setup IS expensive — e.g., generating
			// test data — and you only want to measure the work.
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// The result is intentionally unused. The Go compiler
				// is smart enough that pure side-effect-free calls can
				// be elided, but Sieve allocates memory, which is a
				// side effect — the compiler keeps the call.
				_ = Sieve(n)
			}
		})
	}
}

// fmtName formats a benchmark sub-name like "N=10000". Pulled out as a
// helper to keep the benchmark body readable.
func fmtName(n uint64) string {
	return "N=" + itoa(n)
}

// itoa converts a uint64 to its decimal string. We could use
// strconv.FormatUint, but that's an import we don't otherwise need
// in this file. Keeping the helper local is a tiny pure-style win.
func itoa(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte // uint64 max is 20 digits
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
