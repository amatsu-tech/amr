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
| `amr add [<path>...]` | Add repositories to the workspace (paths are relative to the current directory). `--all` adds every repository found, `-g` puts them into a group |
| `amr remove <path>...` | Remove repositories from the workspace (files on disk are kept) |
| `amr list` | List managed repositories and their groups |
| `amr clone` | Clone repositories that are missing on disk |
| `amr status` | Show branch and changes of repositories |
| `amr pull` | Fast-forward clean repositories from their upstream |
| `amr sync` | Pull the workspace config (if it is a Git repository), clone missing repositories and pull the rest |
| `amr run -- <command>...` | Run a command in repositories |
| `amr root` | Print the workspace root |

Every command except `init` works from any directory inside the workspace.

## Choosing repositories

`list`, `clone`, `status`, `pull`, `sync` and `run` work with every repository by default. Pass paths or groups to work with some of them:

```bash
amr pull frontend services/api        # these repositories
amr pull -g backend                   # repositories in the backend group
amr pull -g backend,web               # in any of these groups (or -g backend -g web)
amr pull infra -g backend             # paths and groups together
amr run -g backend -- go test ./...   # for run, put paths and groups before --
```

Groups are set when adding repositories (`amr add -g backend services/api`) or in the config.

`clone`, `status`, `pull`, `sync` and `run` process up to 8 repositories at the same time. Change it with `--jobs` (`-j`), for example `amr pull -j 16`, or `amr run -j 1 -- <command>` to go one by one.

## Workspace config

A workspace is a folder with `.amr/amr.yaml` at its root. `amr add` fills it in, taking the URL from the repository's `origin`:

```yaml
version: 1
repositories:
  - path: frontend
    url: git@github.com:your-org/frontend.git
  - path: services/api
    url: git@github.com:your-org/api.git
    groups: [backend]
```

A repository without `origin` is saved without `url`, and `amr clone` skips it.

### Sharing a workspace with a team

Keep the config in its own Git repository, for example `your-org/acme-workspace`. `amr init` creates a `.gitignore` that keeps the cloned repositories out of it (unless the folder already has a `.gitignore` or is a Git repository itself):

```gitignore
/*
!/.amr/
!/.gitignore
!/README.md
```

A new teammate then gets the whole layout with:

```bash
git clone git@github.com:your-org/acme-workspace.git acme
cd acme
amr sync
```

Later, `amr sync` keeps everything current: it pulls the config first, so repositories added by others are cloned too.

## Development

Requires Go and [Task](https://taskfile.dev).

```bash
task run -- --help   # run from source
task build           # build binary into build/amr
task check           # format, tidy, build, test
task test            # run tests
task bench           # run benchmarks
```
