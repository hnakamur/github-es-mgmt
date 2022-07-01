VERSION := $(shell git describe --tags HEAD)
DATE := $(shell TZ=UTC git show --quiet --date='format-local:%Y-%m-%dT%H:%M:%SZ' --format="%cd" HEAD)
SRC_FILES = \
  api_error.go \
  client.go \
  cmd/github-es-mgmt/password.go \
  cmd/github-es-mgmt/flag.go \
  cmd/github-es-mgmt/maintenance_enable.go \
  cmd/github-es-mgmt/maintenance_disable.go \
  cmd/github-es-mgmt/settings.go \
  cmd/github-es-mgmt/maintenance.go \
  cmd/github-es-mgmt/settigns_set.go \
  cmd/github-es-mgmt/main.go \
  cmd/github-es-mgmt/certificate_set.go \
  cmd/github-es-mgmt/certificate.go \
  cmd/github-es-mgmt/settings_get.go \
  cmd/github-es-mgmt/maintenance_status.go

all: github-es-mgmt

github-es-mgmt: ${SRC_FILES}
	go build -o github-es-mgmt -trimpath -tags netgo -ldflags "-X main.version=${VERSION} -X main.date=${DATE}" cmd/github-es-mgmt/*.go

clean:
	@rm -f github-es-mgmt
