FROM golang:1.11 as build-env

WORKDIR /go/src/backup-exporter
ADD . /go/src/backup-exporter

RUN go get -d -v ./...
RUN go install

FROM gcr.io/distroless/base
COPY --from=restic/restic:0.9.3 /usr/bin/restic /usr/bin/restic
COPY --from=build-env /go/bin/backup-exporter /
ADD config.toml /
EXPOSE 8080
CMD ["/backup-exporter"]
