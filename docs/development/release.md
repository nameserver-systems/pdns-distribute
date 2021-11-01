# Release
## Build Snapshot Release
To build and test the goreleaser packages, snapshot releases are used. Snapshots use the commit hash as version, instead
of a valid semantic versioned tag. All generated artifacts are stored in a beforehand emptied directory `bin/`.

```bash
make snapshot-release
```

## Build Production Release

!!! important
    Needs environment variable `GITHUB_TOKEN` set for publishing releases and changelogs depending on
    conventional commits.
    

1. Get the latest released tag.
```git
git describe --tags --abbrev=0
```
2. Increment the version, create the tag and push to main repository. Commit tags use the format
[Semantic Versioning](https://semver.org/). The tag description can be the semantic version.
```git
git tag -a v0.0.0 -m "v0.0.0"
git push
```

3. Build release.

!!! caution
    The build and deployment of the production release / binaries is part of the github pipeline.

```bash
make release
```
