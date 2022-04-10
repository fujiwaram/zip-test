
# benchmark

```
$ make bench
go test -bench . -count 1 -benchmem
goos: darwin
goarch: amd64
pkg: fujiwaram/zip-test
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
Benchmark_ArchiveZip_01-12                            68          15057952 ns/op           64362 B/op         59 allocs/op
Benchmark_ArchiveZip_05-12                            15          71131716 ns/op          228603 B/op        131 allocs/op
Benchmark_ArchiveZip_10-12                             8         141481002 ns/op          341201 B/op        216 allocs/op
Benchmark_ArchiveZipByArchiver_01-12               27254             42910 ns/op            3665 B/op         60 allocs/op
Benchmark_ArchiveZipByArchiver_05-12               27565             42960 ns/op            3985 B/op         67 allocs/op
Benchmark_ArchiveZipByArchiver_10-12               27117             43678 ns/op            4362 B/op         73 allocs/op
PASS
ok      fujiwaram/zip-test      11.335s
```