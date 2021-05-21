VERSION := $(shell git describe --contains HEAD)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date --utc +"%Y-%m-%dT%H:%M:%SZ")

all: github-es-mgmt

github-es-mgmt:
	go build -o github-es-mgmt -trimpath -tags netgo -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" cmd/github-es-mgmt/*.go

clean:
	@rm -f github-es-mgmt
