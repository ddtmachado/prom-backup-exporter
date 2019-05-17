package collector

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type backupCollector struct {
	backupRepos     []BackupRepository
	backupSize      *prometheus.Desc
	backupTimestamp *prometheus.Desc
}

// Represents the informations about the snapshot
type BackupSnapshot struct {
	// The snapshot name
	Name string
	// The snapshot creation date and time
	DateString string
	// The snapshot size in bytes
	Size float64
}

// Represents the informations about the backup repositories
type BackupRepository interface {
	// Returns the alias of the repository
	AliasName() string
	// Returns the informations about the latest snapshot of the repository
	LatestSnapshot() (*BackupSnapshot, error)
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewBackupCollector(repos []BackupRepository) *backupCollector {
	labels := []string{"snapshotName", "backupAlias", "creationDate"}

	return &backupCollector{
		backupRepos: repos,
		backupSize: prometheus.NewDesc("backup_size",
			"The size of the backup on the repository",
			labels, nil,
		),
		backupTimestamp: prometheus.NewDesc("backup_timestamp",
			"The minutes elapsed since last backup on the repository",
			labels, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *backupCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.backupSize
	ch <- collector.backupTimestamp
}

//Collect implements required collect function for all promehteus collectors
func (collector *backupCollector) Collect(ch chan<- prometheus.Metric) {

	for _, repo := range collector.backupRepos {
		snapshot, err := repo.LatestSnapshot()
		log.Println("Collect - ")
		log.Println(snapshot)

		if err != nil {
			log.Printf("failed to fetch latest snapshot for %v: %s", repo.AliasName(), err.Error())
			continue
		}

		creationDate, err := time.Parse(time.UnixDate, snapshot.DateString)
		if err != nil {
			log.Printf("failed parse time %v: %s", repo.AliasName(), err.Error())
			continue
		}

		sizeMetric := prometheus.MustNewConstMetric(collector.backupSize, prometheus.GaugeValue, snapshot.Size, snapshot.Name, repo.AliasName(), snapshot.DateString)
		timestampMetric := prometheus.MustNewConstMetric(collector.backupTimestamp, prometheus.GaugeValue, time.Since(creationDate).Minutes(), snapshot.Name, repo.AliasName(), snapshot.DateString)

		ch <- prometheus.NewMetricWithTimestamp(creationDate, sizeMetric)
		ch <- timestampMetric
	}
}
