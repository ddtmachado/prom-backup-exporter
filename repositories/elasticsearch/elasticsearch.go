package elasticsearch

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/ddtmachado/prom-backup-exporter/collector"
)

const snapshotName = "backup-es-latest"

type elasticSearchSnapshot struct {
	Name       string `json:"snapshot"`
	Repository string `json:"repository"`
	Stats      struct {
		TimeInMillis int64   `json:"start_time_in_millis"`
		Size         float64 `json:"total_size_in_bytes"`
	} `json:"stats"`
}

type elasticSearchQuery struct {
	Snapshots []elasticSearchSnapshot `json:"snapshots"`
}

// Represents the informations about the ElasticSearch repository
// used to retrieve informations about the snapshots
type ElasticSearchRepo struct {
	// The repository alias
	Alias,
	// The Elasticsearch URL
	URL,
	// The repository name
	Repo string
}

func OpenRepository(alias, url, repo string) *ElasticSearchRepo {
	return &ElasticSearchRepo{
		Alias: alias,
		URL:   url,
		Repo:  repo,
	}
}

// Returns the alias of the ElasticSearch repository
func (er *ElasticSearchRepo) AliasName() string {
	return er.Alias
}

// Returns the snapshot creation date and time
func (snapshot *elasticSearchSnapshot) DateString() string {
	convertedTime := time.Unix(0, snapshot.Stats.TimeInMillis*int64(time.Millisecond))
	return convertedTime.UTC().Format(time.UnixDate)
}

// Retrieves informations about the latest snapshot of the ElasticSearch repository
func (er *ElasticSearchRepo) LatestSnapshot() (*collector.BackupSnapshot, error) {

	log.Println("Retrieving information about the snapshot - ", snapshotName)

	u, err := url.Parse(er.URL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "_snapshot", er.Repo, snapshotName, "_status")
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http error " + string(resp.StatusCode))
	}

	query := &elasticSearchQuery{}
	err = json.Unmarshal(body, &query)
	if err != nil {
		return nil, err
	}

	if len(query.Snapshots) == 0 {
		return nil, collector.ErrSnapshotNotFound
	}

	return &collector.BackupSnapshot{
		Name:       query.Snapshots[0].Name,
		DateString: query.Snapshots[0].DateString(),
		Size:       query.Snapshots[0].Stats.Size,
	}, nil
}
