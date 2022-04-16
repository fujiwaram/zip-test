package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"

	archiver "github.com/mholt/archiver/v3"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	// CPU pprof
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if err := ArchiveZip("testdata", 1000, "test.zip"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// if err := ArchiveZipByArchiver("testdata", 10, "test_archiver.zip"); err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }

	// Memory pprof
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		// runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	os.Exit(0)
}

func printMemoryStatsHeader() {
	header := []string{"#", "Alloc", "HeapAlloc", "TotalAlloc", "HeapObjects", "Sys", "NumGC"}
	fmt.Println(strings.Join(header, ","))
}

func printMemoryStats(prefix string) {
	// --------------------------------------------------------
	// runtime.MemoryStats() から、現在の割当メモリ量などが取得できる.
	//
	// まず、データの受け皿となる runtime.MemStats を初期化し
	// runtime.ReadMemStats(*runtime.MemStats) を呼び出して
	// 取得する.
	// --------------------------------------------------------
	var (
		ms runtime.MemStats
	)

	runtime.ReadMemStats(&ms)
	data := []string{prefix, toKb(ms.Alloc), toKb(ms.HeapAlloc), toKb(ms.TotalAlloc), toKb(ms.HeapObjects), toKb(ms.Sys), strconv.Itoa(int(ms.NumGC))}
	fmt.Println(strings.Join(data, ","))

	// Alloc は、現在ヒープに割り当てられているメモリ
	// HeapAlloc と同じ.

	// TotalAlloc は、ヒープに割り当てられたメモリ量の累積
	// Allocと違い、こちらは増えていくが減ることはない

	// HeapObjects は、ヒープに割り当てられているオブジェクトの数

	// Sys は、OSから割り当てられたメモリの合計量

	// NumGC は、実施されたGCの回数
}

func toKb(bytes uint64) string {
	return strconv.FormatUint(bytes/1024, 10)
}

func ArchiveZip(inputDir string, maxNum int, outputFilename string) error {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}

	zipFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	printMemoryStatsHeader()

	for i, fi := range files {
		if err := writeOneFile(inputDir, fi, zipWriter); err != nil {
			return err
		}
		if i+1 >= maxNum {
			break
		}
		printMemoryStats(strconv.Itoa(i + 1))
	}
	return nil
}

func writeOneFile(dir string, fi os.DirEntry, zipWriter *zip.Writer) error {
	f, err := os.Open(filepath.Join(dir, fi.Name()))
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := fi.Info()
	if err != nil {
		return err
	}
	header := fileInfoHeader(info)
	header.Method = 8
	w, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}

func fileInfoHeader(fi os.FileInfo) *zip.FileHeader {
	const uint32max = (1 << 32) - 1

	size := fi.Size()
	fh := &zip.FileHeader{
		Name:               fi.Name(),
		UncompressedSize64: uint64(size),
	}
	fh.SetModTime(fi.ModTime())
	fh.SetMode(fi.Mode())
	if fh.UncompressedSize64 > uint32max {
		fh.UncompressedSize = uint32max
	} else {
		fh.UncompressedSize = uint32(fh.UncompressedSize64)
	}
	return fh
}

func ArchiveZipByArchiver(inputDir string, maxNum int, outputFilename string) error {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}
	var filenames []string
	for i, f := range files {
		filenames = append(filenames, filepath.Join(inputDir, f.Name()))
		if i+1 >= maxNum {
			break
		}
	}

	return archiver.Archive(filenames, outputFilename)
}
