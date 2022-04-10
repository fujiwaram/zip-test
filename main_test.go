package main

import "testing"

func Benchmark_ArchiveZip_01(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ArchiveZip("testdata", 1, "test01.zip")
	}
}

func Benchmark_ArchiveZip_05(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ArchiveZip("testdata", 5, "test05.zip")
	}
}

func Benchmark_ArchiveZip_10(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ArchiveZip("testdata", 10, "test10.zip")
	}
}

func Benchmark_ArchiveZipByArchiver_01(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ArchiveZipByArchiver("testdata", 1, "test01_archiver.zip")
	}
}

func Benchmark_ArchiveZipByArchiver_05(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ArchiveZipByArchiver("testdata", 5, "test05_archiver.zip")
	}
}

func Benchmark_ArchiveZipByArchiver_10(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ArchiveZipByArchiver("testdata", 10, "test10_archiver.zip")
	}
}
