
# benchmark

```
$ make bench
go test -bench . -count 1 -benchmem
goos: darwin
goarch: amd64
pkg: fujiwaram/zip-test
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
Benchmark_ArchiveZip_01-12                            76          15193564 ns/op           62126 B/op         64 allocs/op
Benchmark_ArchiveZip_05-12                            15          72920294 ns/op          175479 B/op        151 allocs/op
Benchmark_ArchiveZip_10-12                             7         143531611 ns/op          460789 B/op        260 allocs/op
Benchmark_ArchiveZipByArchiver_01-12                 151           7794627 ns/op           71435 B/op        108 allocs/op
Benchmark_ArchiveZipByArchiver_05-12                 100          36581496 ns/op          213049 B/op        270 allocs/op
Benchmark_ArchiveZipByArchiver_10-12                 100          72606124 ns/op          431900 B/op        468 allocs/op
PASS
ok      fujiwaram/zip-test      19.560s
```