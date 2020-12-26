# Documentation

For creation of this documentation the static site generator [mkdocs](https://www.mkdocs.org/) in
combination with the [mkdocs material theme](https://squidfunk.github.io/mkdocs-material/) and the
[mermaid2 extension](https://github.com/fralau/mkdocs-mermaid2-plugin) is used.
The documentation source are markdown files.
To respect your privacy and being compliant with the GDPR, the generated documentation does not contain external dependencies
like fonts and other CDN resources. Instead of external dependencies all necessary resources are be shipped locally
for example mermaid.js.

## Requirements
* poetry (a python dependency manager)

### Install Environment
1. Python environment with all dependencies must be created/updated by running the command:
```bash
poetry install 
```

### Live Watching Rendered Docs
Documentation can then be viewed at [localhost:8000](http://localhost:8000).

```bash
poetry run mkdocs serve
```

### Build Docs
The static documentation will be generated in the folder `site/`. The content of this folder should be copied to the
documentroot of the documentation webserver.
```bash
poetry run mkdocs build --clean
```

### Update Poetry Environment Dependencies
Updates all dependencies in `pyroject.toml` and `poetry.lock` automatically. If the files were updated, a rerun of
`poetry install` is necessary.

```bash
poetry update
```
