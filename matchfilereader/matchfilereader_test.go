package matchfilereader

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMatchFileReader(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	matches, err := ReadMatchesFile(dir + "/../test/testdata/matchfilereader_testdata.json")
	if err != nil {
		t.Error("Expected nil, got", err)
	}
	if len(matches.Matches) != 1 {
		t.Error("Expected one match in list, got", len(matches.Matches))
	}
}

func TestMatchFileReaderFileNotFound(t *testing.T) {
	matches, err := ReadMatchesFile("/tmp/testtesttest.jsontest")
	if err == nil {
		t.Error("Expected err, got nil")
	}
	if matches != nil {
		t.Error("Expected getting nil, but got an Matches struct")
	}
}
