COMPILEFLAGS = -mod readonly -trimpath -tags netgo,osusergo -ldflags "-s -w -extldflags '-static'"
UPXFLAGS = -7 -t

PATHS = $(shell find . -name "*.go" -not -path "./vendor/*" | xargs -I {} dirname {}  | uniq)
MODULEPATHS = $(shell find . -name "*.go" -not -path "./vendor/*" -not -path "./pkg/*" -not -path "./internal/*" | xargs -I {} dirname {}  | uniq)
MAINMODULE = $(shell go list -m)
MAINDIR = $(shell pwd)
LATESTCOMMIT = $(shell git rev-list --tags --max-count=1)
LATESTGITTAG = $(shell git describe --tags $(LATESTCOMMIT))

all: lint build

build:
	@for f in $(MODULEPATHS); do GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/$$(basename $${f})-$(LATESTGITTAG)-linux-amd64 $(COMPILEFLAGS) $${f}; done

strip-binaries:
	@for f in $(find ./bin/ -name "*-v*") ; do strip -s $${f} ; done


compression:
	@upx ./bin/*

install:
	@for f in $(MODULEPATHS); do go install -i $${f}; done

production-build: COMPILEFLAGS += -a
production-build: mod-prepare lint clean strip-binaries compression

clean:
	go clean -i -testcache ./...
	rm -r -f ./bin/

lint: inspect
inspect: govet golangci-all check-shell-scripts
format: goformat gofix

gofix:
	@for f in $(PATHS); do go fix $${f}; done

goformat:
	@for f in $(PATHS); do go fmt $${f}; done

goimport:
	@goimports -w ./

govet:
	@for f in $(PATHS); do go vet $${f}; done

golangci-all:
	@podman run --rm -t -v ~/.cache/golangci-lint/:/root/.cache -v $(shell pwd):/app -w /app docker.io/golangci/golangci-lint golangci-lint run --fix

check-shell-scripts: generate-goreleaser-config
	@shellcheck ./scripts/package/primary/*.sh
	@shellcheck ./scripts/package/secondary/*.sh

test-coverage:
	go test -race -cover -coverpkg=all ./...

test:
	go test -race -cover ./...

mod-prepare:
	go mod tidy
	go mod download
	go mod verify

download-dep:
	go mod download $(MAINMODULE)

release: generate-goreleaser-config
	goreleaser check
	goreleaser release --clean

snapshot-release: generate-goreleaser-config
	goreleaser check
	goreleaser release --snapshot --clean --skip=publish

generate-goreleaser-config:
	go run cmd/release-config-generator/main.go -ignore release-config-generator,example-microservice

go-update-dependencies:
	go get -u ./...
	go mod tidy

pre-build: generate-goreleaser-config
	goreleaser build --snapshot --clean

.PHONY: build pre-build test download-dep all release snapshot-release go-update-dependencies generate-goreleaser-config
