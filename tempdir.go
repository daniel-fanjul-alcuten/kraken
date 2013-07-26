package kraken

import (
	"io/ioutil"
	"os"
)

type TempDir struct {
	dir   string
	files []string
}

func NewTempDir(dir string) *TempDir {
	return &TempDir{dir, make([]string, 0, 8)}
}

func (td *TempDir) NewFile() (*os.File, error) {
	file, err := ioutil.TempFile(td.dir, "kraken-")
	if err == nil {
		td.files = append(td.files, file.Name())
	}
	return file, err
}

func (td *TempDir) NewDir() (string, error) {
	dir, err := ioutil.TempDir(td.dir, "kraken-")
	if err == nil {
		td.files = append(td.files, dir)
	}
	return dir, err
}

func (td *TempDir) Cleanup() []error {
	errs := make([]error, 0, 2)
	for _, file := range td.files {
		if err := os.RemoveAll(file); err != nil {
			errs = append(errs, err)
		}
	}
	td.files = td.files[:0]
	return errs
}
