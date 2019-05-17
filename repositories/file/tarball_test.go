package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const fileName = "tarball_test_bkp_exporter.tar.gz"
const dirPath = "/tmp/tarball_test"

var fullPath = filepath.Join(dirPath, fileName)

type mockFile struct {
	Name, Date string
	Size       float64
}

func setupTest(t *testing.T) mockFile {
	os.MkdirAll(dirPath, 0777)
	t.Log("Creating file " + fullPath)
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatalln(err)
	}
	file.Close()

	var mock mockFile

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		t.Fatalf("Error retrieving file's informations from path %s - %s ", fullPath, err)
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() == fileName {
			mock.Name = file.Name()
			mock.Date = file.ModTime().UTC().Format(time.UnixDate)
			mock.Size = float64(file.Size())
			break
		}
	}
	return mock
}

func tearDownTest(t *testing.T) {
	t.Log("Removing file " + fullPath)
	os.RemoveAll(dirPath)
}

func TestLatestSnapshot(t *testing.T) {
	t.Log("Testing a success case ")
	createdFile := setupTest(t)
	defer tearDownTest(t)

	repo := OpenRepository("testDir", dirPath, ".tar.gz")

	tarballSnapshot, tarError := repo.LatestSnapshot()

	if tarError != nil {
		t.Errorf("Unexpected error: %s", tarError.Error())
	}

	if tarballSnapshot.Name != fileName {
		t.Errorf("Name - Expected %s but got %s", fileName, tarballSnapshot.Name)
	}

	if tarballSnapshot.DateString != createdFile.Date {
		t.Errorf("Date - Expected %s but got %s", createdFile.Date, tarballSnapshot.DateString)
	}

	if tarballSnapshot.Size != createdFile.Size {
		t.Errorf("Size - Expected %f but got %f", createdFile.Size, tarballSnapshot.Size)
	}

}

func TestLatestSnapshotUnknowFile(t *testing.T) {
	t.Log("Testing an unknow file ")
	repo := OpenRepository("testRepo", "unknow_directory_test_tarball", "tar.gz")
	_, err := repo.LatestSnapshot()
	if err == nil {
		t.Fatalf("Expected an error")
	}
}

func TestLatestSnapshotUnknowExtension(t *testing.T) {
	t.Log("Testing an unknow extension ")
	setupTest(t)
	defer tearDownTest(t)
	repo := OpenRepository("testRepo", dirPath, "*.noext")

	_, tarError := repo.LatestSnapshot()
	if tarError == nil {
		t.Fatalf("Expected error")
	}
}
