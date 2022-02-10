COMPILEFLAGS = -mod readonly -trimpath -tags netgo,osusergo -ldflags "-s -w -extldflags '-static'"
UPXFLAGS = -7 -t

PATHS = $(shell find . -name "*.go" -not -path "./vendor/*" | xargs -I {} dirname {}  | uniq)
MODULEPATHS = $(shell find . -name "*.go" -not -path "./vendor/*" -not -path "./pkg/*" -not -path "./internal/*" | xargs -I {} dirname {}  | uniq)
MAINMODULE = $(shell go list -m)
MAINDIR = $(shell pwd)
LATESTCOMMIT = $(shell git rev-list --tags --max-count=1)
LATESTGITTAG = $(shell git describe --tags $(LATESTCOMMIT))

CONTAINERFILES="./build/ci/powerdns" "./build/package/pdns-primary"

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
	golangci-lint run --fix ./...

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

build-images:
	@for f in $(CONTAINERFILES); do	podman build -t $$(basename $${f}) $${f}; done

cleanup:
	podman image prune -f -a
	podman container prune -f

start-dev-dependencies:
	@for f in consul nats pdns-primary pdns-secondary ; do podman start $$(basename $${f}) ; podman port $$(basename $${f}) ; done

start-create-dev-dependencies: build-images
	podman run -d -p 8300:8300 -p 8500:8500 -p 8600:8600 --name consul consul
	podman port consul
	podman run -d -p 4222:4222 -p 6222:6222 -p 8222:8222 --name nats nats -m 8222
	podman port nats
	podman run -d -p 8081:8081 -p 5353:53 --name pdns-primary pdns-primary
	podman port pdns-primary
	podman run -d -p 18081:8081 -p 53535:53 --name pdns-secondary pdns-primary
	podman port pdns-secondary

stop-dev-containers:
	podman stop consul
	podman stop nats
	podman stop pdns-primary
	podman stop pdns-secondary

stop-all-containers:
	podman stop -a

release: generate-goreleaser-config
	goreleaser check
	goreleaser release --rm-dist

snapshot-release: generate-goreleaser-config
	goreleaser check
	goreleaser release --snapshot --rm-dist --skip-publish

generate-goreleaser-config:
	go run cmd/release-config-generator/main.go -ignore release-config-generator,example-microservice

go-update-dependencies:
	go get -u ./...
	go mod tidy

pre-build: generate-goreleaser-config
	goreleaser build --snapshot --rm-dist

.PHONY: build pre-build test download-dep all release snapshot-release go-update-dependencies generate-goreleaser-config
