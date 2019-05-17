# Builds `backup-exporter` to the current location.
build:
	go build -i -v


# Builds and installs `backup-exporter` to the currently
# configured $GOPATH location for binaries.
install:
	go install -v


# Formats all the non-vendors source code.
fmt:
	go fmt ./...


# Builds the Docker image.
image:
	./build-image.sh

# Runs all the tests specified under the root directory
# (excluding vendor).
#
test:
	mkdir -p ./test-results && \
	go test -failfast -short -v ./... 2>&1 | go-junit-report > test-results/TEST-report.xml

run-local:
	docker-compose --project-name exporter up -d
