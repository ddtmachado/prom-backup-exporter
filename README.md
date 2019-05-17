# Overview
Backup Exporter is responsible for collecting and exporting metrics from the latest backups of pre-configured repositories.

Currently it supports collecting from Elasticsearch, Tarball and Restic repositories.

## Dependencies

[Go 1.11] (https://golang.org/doc/install)

## Configuration file

The config.toml is the default configuration file and must be created in the /etc/backup-exporter directory.

[It's possible to use another configuration file using the flag --config.](#using-flags)

This file must be configured with the follow informations:

```
port = 8080       ## Port where Prometheus is running
path = '/metrics' ## Path where Prometheus collects metrics

## Repositories - must be restic, tarball or elasticsearch

## More than one entry for the same kind of repository -> [[repository name]]

[[restic]]                ## Restic repository configuration
  alias = 'tagExample1'   ## Tag used on the snapshot creation
  path = 'tmp/restic'     ## The location of the restic repository
  password = 'pass'       ## The password to restic repository access

[[restic]]
  alias = 'tagExample2'
  path = 'repository/restic'
  password = 'anotherpass'

## Only one entry for the repository -> [repository name]

[tarball]                 ## Tarball repository configuration
  alias = 'wdBackups'     ## The repository alias
  path = '/backups'       ## The location of tarball repository
  extension = '.tar.gz'   ## The extension file to be filtered

[elasticsearch]                   ## Elasticsearch repository configuration
  alias = 'elasticsearch-shared'  ## The repository alias
  url = 'http://localhost:9200/'  ## The Elasticsearch URL
  repo = 'es_repo'                ## The repository name
```

## Running Backup Exporter

On the root directory type:

```sh
go run main.go
```

### Using flags

It's possible to run the Backup Exporter application with the following flags, which will override the config file if present:

- --config  - The full path of the configuration file to be used
- --port    - The port where Prometheus is running
- --path    - The path where Prometheus collects metrics

Example:

```
go run main.go --config /etc/backup-exporter/config/dev-config.toml --port 9090 --path /new-metrics
```

## Running the tests

To run all the unit tests, just run the following in the root directory of this repository:

```
go test ./...
```

To run the unit tests for a specific package, just run `go test` from the package folder
