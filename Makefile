VERSION := $(shell git describe --tags HEAD)
DATE := $(shell TZ=UTC git show --quiet --date='format-local:%Y-%m-%dT%H:%M:%SZ' --format="%cd" HEAD)

all: github-es-mgmt

github-es-mgmt:
	go build -o github-es-mgmt -trimpath -tags netgo -ldflags "-X main.version=${VERSION} -X main.date=${DATE}" cmd/github-es-mgmt/*.go

clean:
	@rm -f github-es-mgmt
