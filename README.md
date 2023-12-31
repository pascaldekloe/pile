# Pile

Pile provides sorted data structures for the Go programming language.

The Map operations are Find, Insert, Update and Put, plus Swap from Iterator.
Delete is not implemented (yet). Iterator instantiation with At, Least or Most
is lightweight—no memory alloctaion.

This is free and unencumbered software released into the
[public domain](https://creativecommons.org/publicdomain/zero/1.0).

[![API Documentation](https://godoc.org/github.com/pascaldekloe/pile?status.svg)](https://godoc.org/github.com/pascaldekloe/pile)
[![Build Status](https://github.com/pascaldekloe/pile/actions/workflows/go.yml/badge.svg)](https://github.com/pascaldekloe/pile/actions/workflows/go.yml)


# Scale

The Map works well with large Value types and no pointer(s). Data in the Map
releases the garbage collector with fewer objects to care about. Go's native
(hash) map can't cope with embed values at scale as resize events eventually
become too big.


# Performance

The implementation shines when reads or writes operate on a small or predictable
key range. Without the CPU cache in our favour, operation slows down to a little
over one million per second.

```
goos: darwin
goarch: arm64
pkg: github.com/pascaldekloe/pile
BenchmarkMapFind/Sequential/1Ki-8     	136814722	         8.288 ns/op
BenchmarkMapFind/Sequential/1Mi-8     	39862473	        30.96 ns/op
BenchmarkMapFind/Sequential/64Mi-8    	35958916	        33.09 ns/op
BenchmarkMapFind/Random/1Ki-8         	79347367	        13.21 ns/op
BenchmarkMapFind/Random/1Mi-8         	 3581359	       333.5 ns/op
BenchmarkMapFind/Random/64Mi-8        	 1595295	       778.6 ns/op
BenchmarkInsert/Append/map-8          	28757638	        41.31 ns/op
BenchmarkInsert/Append/set-8          	29841837	        39.63 ns/op
BenchmarkInsert/Prepend/map-8         	28124312	        46.20 ns/op
BenchmarkInsert/Prepend/set-8         	29796175	        42.58 ns/op
BenchmarkInsert/Random/map-8          	 4317328	       493.5 ns/op
BenchmarkInsert/Random/set-8          	 5090876	       470.5 ns/op
PASS
ok  	github.com/pascaldekloe/pile	111.711s
```

Pile has its own strengths and weaknesses when
[compared to others](https://github.com/tidwall/btree-benchmark/pull/4).
