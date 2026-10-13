# amr

**amr** (Amatsu multi-repo) is a CLI for managing many Git repositories as one workspace.

Clone every repository with one command, see the status of all of them at once, pull everything, or run a command in each repository.

> Work in progress. Commands and the config format may change before v1.

## Quick start

```bash
cd ~/projects/my-company
amr init        # create a workspace here
amr scan        # find Git repositories, including nested ones like services/api
amr add --all   # add every repository found
amr status      # branch and changes of every repository
```

On another machine, copy `.amr/amr.yaml` into an empty folder and run `amr clone` to get the same layout.

## Commands

| Command | What it does |
| --- | --- |
| `amr init [directory]` | Create a workspace in the current or given directory |
| `amr scan` | Find Git repositories in the workspace and show which are managed |
| `amr add [<path>...]` | Add repositories to the workspace (paths are relative to the current directory). `--all` adds every repository found |
| `amr remove <path>...` | Remove repositories from the workspace (files on disk are kept) |
| `amr list` | List managed repositories |
| `amr clone` | Clone every managed repository that is missing on disk |
| `amr status` | Show branch and changes of every managed repository |
| `amr pull` | Fast-forward every clean repository from its upstream |
| `amr run -- <command>...` | Run a command in every managed repository |
| `amr root` | Print the workspace root |

Every command except `init` works from any directory inside the workspace.

`clone`, `status`, `pull` and `run` process up to 8 repositories at the same time. Change it with `--jobs` (`-j`), for example `amr pull -j 16`, or `amr run -j 1 -- <command>` to go one by one.

## Workspace config

A workspace is a folder with `.amr/amr.yaml` at its root. `amr add` fills it in, taking the URL from the repository's `origin`:

```yaml
version: 1
repositories:
  - path: frontend
    url: git@github.com:your-org/frontend.git
  - path: services/api
    url: git@github.com:your-org/api.git
```

A repository without `origin` is saved without `url`, and `amr clone` skips it.

## Development

Requires Go and [Task](https://taskfile.dev).

```bash
task run -- --help   # run from source
task build           # build binary into build/amr
task check           # format, tidy, build
```
