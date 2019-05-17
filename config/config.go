package config

import (
	"github.com/ddtmachado/prom-backup-exporter/collector"
	"github.com/ddtmachado/prom-backup-exporter/repositories/elasticsearch"
	"github.com/ddtmachado/prom-backup-exporter/repositories/file"
	"github.com/ddtmachado/prom-backup-exporter/repositories/restic"
)

// Config represents the configuration struct of the service.
type Config struct {
	//HTTP port that will be used
	Port string
	//HTTP path to export metrics, defaults to "/metrics"
	Path string

	ResticRepos        []*restic.ResticRepository         `mapstructure:"restic"`
	ElasticSearchRepos []*elasticsearch.ElasticSearchRepo `mapstructure:"elasticsearch"`
	TarballRepos       []*file.TarballRepo                `mapstructure:"tarball"`
}

// Repos returns a concatenated list of all repositories found
// in the config file.
func (c *Config) Repos() []collector.BackupRepository {
	var repos []collector.BackupRepository
	for _, repo := range c.ResticRepos {
		repos = append(repos, repo)
	}
	for _, repo := range c.ElasticSearchRepos {
		repos = append(repos, repo)
	}
	for _, repo := range c.TarballRepos {
		repos = append(repos, repo)
	}
	return repos
}
