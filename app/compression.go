package app

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func countRec(source string) uint64 {
	var count uint64
	f, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	stat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		return 1
	}
	files, err := ioutil.ReadDir(source)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			count += countRec(source + "/" + file.Name())
		}
		count++
	}
	return count
}

func CreateArchive(source, target, name string, pchan chan float64) (*os.File, error) {
	var file float64
	count := countRec(source) + 1
	log.Println("Files:", count)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", name))
	tarfile, err := os.Create(target)
	if err != nil {
		return nil, err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil, err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return tarfile, filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			file++
			percentage := float64(100.0 * (float32(file) / float32(count)))
			log.Println("Progress:", percentage, "%")
			if pchan != nil {
				pchan <- percentage
			}
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}
