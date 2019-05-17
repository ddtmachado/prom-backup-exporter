package file

import (
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/ddtmachado/prom-backup-exporter/collector"
)

// Represents the informations about the tarball respository
// used to retrieve informations about the snapshots
type TarballRepo struct {
	// The repository alias
	Alias,
	// The path to the repository
	Path,
	// The extension of the files to be checked
	Extension string
}

func addDotToFileExtension(fileExtension string) string {
	if fileExtension != "" && !strings.HasPrefix(fileExtension, ".") {
		fileExtension = "." + fileExtension
	}
	return fileExtension
}

func OpenRepository(alias, path, extension string) *TarballRepo {
	return &TarballRepo{
		Alias:     alias,
		Path:      path,
		Extension: extension,
	}
}

// Returns the repository alias
func (t *TarballRepo) AliasName() string {
	return t.Alias
}

// Retrieves informations about the latest snapshot of the Tarball repository
func (t *TarballRepo) LatestSnapshot() (*collector.BackupSnapshot, error) {

	log.Println("Retrieving snapshot informations from path - ", t.Path)

	files, err := ioutil.ReadDir(t.Path)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	fileExtension := addDotToFileExtension(t.Extension)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), fileExtension) {
			return &collector.BackupSnapshot{
				Name:       file.Name(),
				DateString: file.ModTime().UTC().Format(time.UnixDate),
				Size:       float64(file.Size()),
			}, nil
		}
	}
	return nil, collector.ErrSnapshotNotFound
}
