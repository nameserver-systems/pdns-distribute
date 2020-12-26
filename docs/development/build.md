# Build

## Requirements
* Go > 1.15.X
* goreleaser
* upx
* golangci-lint
* podman

## Building Binaries
Binaries are build like snapshot releases (equal to [Snapshot Release](release.md#build-snapshot-release)).
The binaries are in the directory `bin/`.

```bash
make snapshot-release
```
