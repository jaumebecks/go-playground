# My personal Go Playground

## File Generator

### Problem statement

Given a dummy feed data source, create a CSV from its data iteration

### Sources

`file-generator-db` sqlite3 database contains a table `main.feed` which includes

| id_item                 | id_offer | price | title  | brand  | category | in_promo |
|-------------------------|----------|-------|--------|--------|----------|----------|
| int (PK auto-increment) | int      | float | string | string | string   | boolean  |


### Benchmark

```shell
go test -bench=. file-generator -benchtime=10x
```

#### Results

When dealing with concurrency, file generation time is decreased by 10x,
directly proportional to the amount of goroutines provided

```shell
$ go test -bench=. file-generator -benchtime=10x

goos: linux
goarch: amd64
pkg: file-generator
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz

BenchmarkGenerateFeedSequentially-12    	2023/02/16 15:55:30 GenerateSpecificFeed1Sequentially, execution time 10.319924329s
      10	10312030383 ns/op

BenchmarkGenerateConcurrently-12        	2023/02/16 15:57:05 GenerateSpecificFeed1Concurrently, execution time 1.042788643s
      10	1054748184 ns/op

PASS
ok  	file-generator	125.099s
```
