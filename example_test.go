package primesingo_test

import (
	"fmt"

	"github.com/ezra-gocci/primesingo"
)

// ExampleSieve demonstrates basic usage. The "// Output" comment makes
// this example executable as a test — go test runs the function and
// verifies stdout matches.
func ExampleSieve() {
	ps := primesingo.Sieve(20)
	fmt.Println(ps)
	// Output: [2 3 5 7 11 13 17 19]
}
