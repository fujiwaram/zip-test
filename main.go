package main

import (
	"archive/zip"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"

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

	if err := ArchiveZip("testdata", 10, "test.zip"); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := ArchiveZipByArchiver("testdata", 10, "test_archiver.zip"); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// Memory pprof
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	os.Exit(0)
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

	for i, fi := range files {
		if err := writeOneFile(inputDir, fi, zipWriter); err != nil {
			return err
		}
		if i+1 >= maxNum {
			break
		}
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
