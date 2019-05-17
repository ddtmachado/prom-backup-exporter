package restic

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/ddtmachado/prom-backup-exporter/collector"
)

var execCommand = exec.Command

// Represents the informations about the restic respository
// used to retrieve informations about the snapshots
type ResticRepository struct {
	// The repository alias - used as the restic backup tag
	Alias,
	// The path to the repository
	Path,
	// The required password to open the restic repository
	Password string
}

type resticSnapshot struct {
	Time        string
	Tags, Paths []string
	Id          string `json:"short_id"`
}

type resticSnapshotStats struct {
	TotalSize      float64 `json:"total_size"`
	TotalBlobCount float64 `json:"total_blob_count"`
}

func (r *ResticRepository) open() error {
	cmd := execCommand("restic", "check")
	cmd.Env = r.environmentVariables()
	_, err := cmd.Output()
	return err
}

func (r *ResticRepository) environmentVariables() []string {
	return append(os.Environ(),
		"RESTIC_PASSWORD="+r.Password,
		"RESTIC_REPOSITORY="+r.Path,
	)
}

func (r *ResticRepository) exec(args ...string) ([]byte, error) {
	cmd := execCommand("restic", append(args, "--json")...)
	cmd.Env = r.environmentVariables()
	out, err := cmd.CombinedOutput()
	log.Printf("restic process output: %s", out)
	return out, err
}

func (r *ResticRepository) latestSnapshotForTag() (*resticSnapshot, error) {
	var snapshot []resticSnapshot
	out, err := r.exec("snapshots", "--last", "--tag", r.Alias)
	if err != nil {
		log.Println(err)
		return &resticSnapshot{}, err
	}

	err = json.Unmarshal(out, &snapshot)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	if len(snapshot) == 0 {
		return &resticSnapshot{}, collector.ErrSnapshotNotFound
	}
	return &snapshot[0], nil
}

func (snapshot *resticSnapshot) creationDateString() string {
	creationDate, err := time.Parse(time.RFC3339, snapshot.Time)
	if err != nil {
		log.Println(err)
		return ""
	}
	return creationDate.UTC().Format(time.UnixDate)
}

func (r *ResticRepository) snapshotSize(snapshot *resticSnapshot) float64 {
	out, err := r.exec("stats", snapshot.Id, "--mode", "raw-data")
	if err != nil {
		log.Println(err)
		return 0
	}

	var snapshotStat resticSnapshotStats
	err = json.Unmarshal(out, &snapshotStat)
	if err != nil {
		fmt.Println("error:", err)
	}

	return snapshotStat.TotalSize
}

// Returns the backup tag name
func (r *ResticRepository) AliasName() string {
	return r.Alias
}

// Retrieves informations about the latest snapshot of the Restic repository
func (r *ResticRepository) LatestSnapshot() (*collector.BackupSnapshot, error) {
	resticSnapshot, err := r.latestSnapshotForTag()
	if err != nil {
		return nil, err
	}
	return &collector.BackupSnapshot{
		Name:       resticSnapshot.Id,
		DateString: resticSnapshot.creationDateString(),
		Size:       r.snapshotSize(resticSnapshot),
	}, nil
}
