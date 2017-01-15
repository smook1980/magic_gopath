package main

// Magic Go Path
//
// Inspired by http://www.jtolds.com/writing/2017/01/magic-gopath/

import (
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	dirName, err := os.Getwd()

	if err != nil {
		log.Fatalln("Unable to get current working directory: ", err)
		os.Exit(-2)
	}

	fmt.Println(findGoPath(dirName))
	os.Exit(0)
}

func findGoPath(dirpath string) string {
	if gopath := ReadDotFilePath(dirpath); gopath != "" {
		return gopath
	}

	if IsBaseGoPath(dirpath) {
		return dirpath
	}

	if dirpath == "/" {
		log.Fatalln("Unable to locate a suitable GOPATH.")
		os.Exit(-1)
	}

	return findGoPath(path.Dir(dirpath))
}

func IsBaseGoPath(dirpath string) bool {
	f, err := os.Open(dirpath)

	if err != nil {
		log.Fatalln("Could not stat path ", dirpath)
		os.Exit(-3)
	}

	defer f.Close()

	dirs, err := f.Readdirnames(0)

	if err != nil {
		log.Fatalln("Failed to list dir ", dirpath)
		os.Exit(-4)
	}

	matchedPaths := 0
	gopathDirs := []string{"src", "pkg", "bin"}

	for _, d := range dirs {
		for _, gd := range gopathDirs {
			if d == gd {
				matchedPaths += 1
			}
		}
	}

	return matchedPaths > 1
}

// ReadDotFilePath attempts to read .gopath in the given fpath
//
// Return the contents of .gopath if the file is readable, otherwise
// returns nil if file does not exist or can't be read due to error.
// All errors are logged to STDERR
func ReadDotFilePath(fpath string) string {
	inf, err := os.Stat(fpath)

	if os.IsNotExist(err) {
		return ""
	} else if os.IsPermission(err) {
		log.Println("Permission denied accessing ", fpath)
		return ""
	} else if err != nil {
		log.Println("Unknown error accessing ", fpath, ": ", err)
		return ""
	}

	if inf.Mode().IsRegular() {
		f, err := os.Open(fpath)

		if err != nil {
			log.Println("Unknown error reading ", fpath)
			return ""
		}

		defer f.Close()

		byteBuff := make([]byte, int(inf.Size()))
		_, err = f.Read(byteBuff)

		if err != nil {
			log.Println("Unknown error reading ", fpath)
			return ""
		}

		return string(byteBuff)
	}

	return ""
}
