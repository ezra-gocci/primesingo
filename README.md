# primesingo

A Go library for enumerating prime numbers using a parallel segmented Sieve of Eratosthenes.

## Status

Educational project, in active development.

## Install

    go get github.com/ezra-gocci/primesingo

## Usage

    ps := primesingo.Sieve(100)
    fmt.Println(ps) // [2 3 5 7 ... 97]

## Algorithm

See package documentation for design notes.

## License

MIT

## Learning Q&A

This project is built as a learning exercise. The questions below come from
comprehension checks during implementation — each one captures a decision
point or subtlety worth understanding.

### Algorithm

**Q: Why use `p*p <= n` instead of `p <= sqrt(n)` as the loop condition?**

A: Two reasons. First, integer multiplication is faster than calling
`math.Sqrt`. Second and more important: `math.Sqrt` returns `float64`,
which can't exactly represent most integers above 2⁵³. A rounding error
of one in either direction can either waste a loop iteration or — worse —
stop one short and miss a divisor, producing wrong results. Integer
arithmetic `p*p <= n` has no floating-point in it at all and is exact.

**Q: Why do we only sieve primes up to √n?**

A: Any composite c ≤ n factors as c = a·b with a ≤ b. Then a² ≤ a·b ≤ n,
so a ≤ √n. The smallest prime factor of any composite ≤ n is therefore
≤ √n. Once we've marked multiples of every prime up to √n, every
composite ≤ n has been marked. Primes greater than √n have no unmarked
work to do — their multiples are all > n.

**Q: Why start marking multiples at p² instead of 2p?**

A: Any multiple of p smaller than p² has the form p·k for some k < p.
That k has a prime factor q ≤ k < p, so when we processed q earlier, we
already marked p·k as a multiple of q. Starting at p² skips redundant
writes.

**Q: Why is the time complexity O(n log log n) and not just O(n log n)?**

A: Total work is N · Σ(1/p) for primes p ≤ √N. By Mertens' theorem,
Σ(1/p) over primes ≤ N grows like ln(ln(N)), not ln(N). The outer log
comes from harmonic-series-like density of integers; the inner log comes
from the Prime Number Theorem (primes thin out as ~1/ln(n)). Composing
them gives log log.

**Q: Why use a separate collection loop instead of collecting primes
during the marking phase?**

A: The marking loop only iterates p up to √n. Primes greater than √n —
the majority of primes — are never visited there; they live as unmarked
indices that only the collection loop sees. Collecting during marking
would require a second pass anyway, splitting one clean loop into two
fragments. The two-phase structure (mark all, then collect all) is
cleaner and equally efficient.

### Go: Types and Memory

**Q: Why `uint64` for the parameter instead of `int`?**

A: Primes are non-negative, so an unsigned type encodes the constraint in
the type system. Also: when we later use bit arithmetic (`i >> 6`,
`i & 63`), unsigned types behave predictably; signed right-shift
sign-extends, which is not what bit-packing wants.

**Q: Why `make([]bool, n+1)` and not `make([]bool, n)`?**

A: We want to index the slice directly by the number being tested, so
asking "is `n` composite?" is `composite[n]`. For `composite[n]` to be a
valid index, the slice must have at least `n+1` elements (indices 0
through n). The slots for 0 and 1 are unused but the readability win
from "index equals the number" justifies the 2 bytes.

**Q: Why name the slice `composite` instead of `isPrime` or `sieve`?**

A: The marking loop only ever *sets* entries to `true` (never `false`).
Naming the slice for what `true` means makes the code read naturally —
"if composite[m] is true, skip; else mark composite[m] = true." Inverting
to `isPrime` would require initializing every entry to `true` first (an
O(n) pass), or flipping the meaning of the marking. The zero value of
`bool` is `false`, which means our chosen direction gives us
initialization for free.

**Q: Could `p*p` overflow uint64?**

A: For practical inputs, no. The loop condition `p*p <= n` means p stops
growing once p² exceeds n. For n < 2⁶⁴ (always true for uint64), p never
exceeds ~2³² before the loop terminates, so p*p stays below uint64 max.
The real overflow risk is `m += p` in the inner loop when n is near
MaxUint64 — m can wrap to a small value and produce an infinite loop.
For workstation-scale n this is purely theoretical.

### Go: Testing

**Q: Why does the example test live in `package primesingo_test` while
the unit tests live in `package primesingo`?**

A: The `_test` suffix is the only place in Go where two packages share a
directory. Tests in `package primesingo` have white-box access to
unexported identifiers — useful for testing internal helpers. The
example file uses the black-box `primesingo_test` package, so it must
import the public API exactly as a real consumer would. This keeps
examples honest: they cannot accidentally depend on internal state.

**Q: What's the difference between `t.Errorf` and `t.Fatalf`?**

A: `Errorf` records a failure and continues running the test. `Fatalf`
records and stops the current test (or subtest) immediately. Use
`Fatalf` when later assertions depend on earlier success (e.g., don't
index into a slice that you just confirmed is nil). Use `Errorf` for
independent assertions so you see all failures at once.

**Q: Why use `reflect.DeepEqual` to compare slices instead of `==`?**

A: Slice comparison with `==` is a compile error in Go (slices are not
comparable). `reflect.DeepEqual` does element-wise comparison and
handles nil correctly. Note: `DeepEqual(nil, []uint64{})` returns
`false` — nil and empty-but-allocated slices are not equal under
DeepEqual. As of Go 1.21+, `slices.Equal` is faster (no reflection) and
treats nil and empty as equal. Pick the one whose semantics match what
you're testing.

**Q: What does `testing.Short()` do?**

A: Returns `true` if `go test` was run with `-short`. Tests can consult
it via `if testing.Short() { t.Skip(...) }` to skip expensive cases
during fast-feedback runs. Convention: developers run `-short` locally
during the edit-test loop; CI runs the full suite without `-short`.

**Q: Why use a reference implementation (trial division) inside the
test file?**

A: It's a *test oracle* — a slow, obviously-correct implementation used
to validate the fast, clever one. Whenever you write a fast algorithm,
write a slow reference alongside it in the test file. Living in the
test file keeps it correctly scoped (not part of the public API) and
out of production code.

**Q: How does `-run` actually filter tests?**

A: `-run` is an unanchored substring regex match against function names.
It filters `TestXxx` and `ExampleXxx` functions (but not `BenchmarkXxx`,
which has its own `-bench` flag). `-run TestSieve` matches all
`TestSieve_*` functions; `-run Sieve` also matches `ExampleSieve`. Use
`^...$` anchors for strict matching.

### Go: Tooling

**Q: Why does `go doc .` show only signatures by default?**

A: It's optimized for "remind me what's in this package" lookups, where
you already know what symbols do. Use `go doc -all .` for full
documentation including the package overview and all doc comments. For
browsing, run `godoc -http=:6060` (or `pkgsite -http=:8080`) and view
the rendered docs in a browser.

**Q: Why is package-level documentation conventionally in `doc.go`?**

A: When a package has substantial documentation, it lives in its own
file rather than cluttering the top of an implementation file. The
package comment must be immediately above the `package` declaration —
no blank line between — for godoc to associate them.

**Q: Should the GitHub repo name match the Go package name?**

A: Not necessarily. The module path (`go.mod`) must match the GitHub
URL because Go uses that path to fetch the module. But the package
name (declared in `.go` files) is independent — callers reference the
package by its declared name, not the URL. Repo names can have hyphens
or underscores; Go package names should be short, lowercase, and
without underscores.