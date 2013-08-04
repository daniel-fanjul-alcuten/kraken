package kraken

import (
	"os"
	"testing"
)

func TestTempDirNewFile(t *testing.T) {

	td := NewTempDir("")
	if td == nil {
		t.Fatal()
	}
	defer td.Cleanup()
	if td.dir != "" {
		t.Error(td.dir)
	}

	file, err := td.NewFile()
	if err != nil {
		t.Error(err)
	}
	if file == nil {
		t.Fatal(err)
	}

	name := file.Name()
	if err := file.Close(); err != nil {
		t.Error(err)
	}
	if _, err := os.Stat(name); err != nil {
		t.Error(err)
	}

	if len(td.files) != 1 {
		t.Error(len(td.files))
	}
	if td.files[0] != name {
		t.Error(td.files[0])
	}

	for _, err := range td.Cleanup() {
		t.Error(err)
	}
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Error(err)
	}
	if len(td.files) != 0 {
		t.Error(len(td.files))
	}
}

func TestTempDirNewDir(t *testing.T) {

	td := NewTempDir("")
	if td == nil {
		t.Fatal()
	}
	defer td.Cleanup()
	if td.dir != "" {
		t.Error(td.dir)
	}

	name, err := td.NewDir()
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(name); err != nil {
		t.Error(err)
	}

	if len(td.files) != 1 {
		t.Error(len(td.files))
	}
	if td.files[0] != name {
		t.Error(td.files[0])
	}

	for _, err := range td.Cleanup() {
		t.Error(err)
	}
	if _, err := os.Stat(name); !os.IsNotExist(err) {
		t.Error(err)
	}
	if len(td.files) != 0 {
		t.Error(len(td.files))
	}
}
