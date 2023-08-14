# My personal Go Playground

## Struct assignment

### Problem statement

Given a dummy string slice data source, create a slice of structs, both
sequentially, and concurrently

### Sources

In [main.go](main.go), `GenerateInput()` will generate a slice as large as
the alphabet (_27 in our case_) multiplied by the const `NumElements`, for example

> const NumElements = 10000, will define a slice of 270.000 strings
> ["a", "b", "c", "d", ..., "z"] * 10000 times

### Solution approach

Given a result set from the provided data source, the idea is to generate a slice
of output structs.

The approach consist on having one method for each methodology (sequential/parallel)
which will share common functionalities such as creating the output struct.

#### Sequential code
The idea is to iterate over the total amount of input, and directly create
data structures

#### Parallel code
The idea is to split the total amount of input rows into smaller chunks,
spin goroutines for each chunk, which will generate their corresponding structs.

This will have 2 variations, 1 shared map for storing data from goroutines
(the Map approach), and on the other side 1 channel for communicating each goroutine
output (the Channel approach)

This process will get benefit of `sync.WaitGroup`, builtin functionality, to control
the amount of workers that will generate the structs

### Benchmark

```shell
go test -bench=.
```

#### Results

When dealing with parallelism, struct assignment is not getting benefits from
parallelism. Sequential assignment is being faster when dealing with any amount
of strings (Tested with 27, 270.000 & 2.700.000 items, with same result)

```shell
$ go test -bench=.

goos: linux
goarch: amd64
pkg: struct-assignment
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkGenerateFeedSequentially-12           	     357	   3221287 ns/op
BenchmarkGenerateConcurrentlyWithMap-12        	     103	  11120158 ns/op
BenchmarkGenerateConcurrentlyWithChannel-12    	      55	  21139453 ns/op
PASS
ok  	struct-assignment	4.818s
```
