# primes

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