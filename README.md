# amr

**amr** (Amatsu multi-repo) is a CLI for managing many Git repositories as one workspace.

Clone every repository with one command, see the status of all of them at once, pull everything, or run a command in each repository.

> Work in progress. Commands and the config format may change before v1.

## Workspace config

A workspace is a folder with an `amr.yaml` file at its root:

```yaml
version: 1
repositories:
  - path: frontend
    url: git@github.com:your-org/frontend.git
  - path: backend/catalog
    url: git@github.com:your-org/catalog.git
```
