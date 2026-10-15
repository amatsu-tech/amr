# amr

English | [Русский](README.ru.md)

**amr** (Amatsu multi-repo) is a command-line tool that helps you work with many Git repositories as one workspace.

Clone every repository with one command, see the state of each one, update everything at once, or run a command in all repositories.

> Work in progress. Commands and the config format may change before version 1.0.

## Quick start

```bash
cd ~/projects/my-company
amr init        # create a workspace in this folder
amr scan        # find Git repositories, including nested folders like services/api
amr add --all   # add every repository found
amr status      # branch and changes of every repository
```

To share the workspace with your team and get the same folder layout on any machine, see [Workspace for a team](#workspace-for-a-team).

## What a workspace is

A workspace is a folder that holds your repositories and has a service folder `.amr` with the config at its root. For example:

```
my-company/
├── .amr/amr.yaml      ← list of repositories and their URLs
├── frontend/
├── infra/
└── services/
    ├── api/
    └── billing/
```

Every command except `init` can be run from any folder inside the workspace: `amr` walks up the tree and finds `.amr` by itself.

## Commands

| Command | What it does |
| --- | --- |
| `amr init [directory]` | Create a workspace in the current or given folder |
| `amr scan` | Find Git repositories in the workspace and show which are already added |
| `amr add [<path>...]` | Add repositories to the workspace (paths are relative to the current folder). `--all` adds every repository found, `-g` puts them into a group |
| `amr remove <path>...` | Remove repositories from the workspace (files on disk are kept) |
| `amr list` | Show added repositories and their groups |
| `amr clone` | Clone repositories that are missing on disk |
| `amr status` | Show branch and changes of repositories |
| `amr pull` | Pull changes into repositories without local changes (fast-forward only) |
| `amr sync` | Update the workspace config (if it is a Git repository), clone what is missing and pull the rest |
| `amr run -- <command>...` | Run a command in repositories |
| `amr stash` | Stash changes in repositories under one name, including new files. See [Stashing changes](#stashing-changes) |
| `amr root` | Print the path to the workspace root |

Details for each command: `amr <command> --help`.

## Choosing repositories

`list`, `clone`, `status`, `pull`, `sync`, `run` and `stash` work with every repository by default. To work with only some of them, pass paths or groups:

```bash
amr pull frontend services/api        # only these repositories
amr pull -g backend                   # repositories in the backend group
amr pull -g backend,web               # in any of these groups (or -g backend -g web)
amr pull infra -g backend             # a path and a group together
amr pull .                            # the repository you are in, even from its subfolder
amr run -g backend -- go test ./...   # for run, put paths and groups before --
```

When you pass both paths and groups, amr takes every matching repository: those given by path and those in the groups. The order is the same as in the config, without duplicates.

If there is no such path or group in the config, for example because of a typo, the command does nothing and shows an error right away:

```
$ amr pull -g bakend
Error: no repositories in group: bakend
```

Groups are set when adding repositories (`amr add -g backend services/api`) or directly in the config.

## Stashing changes

`amr stash` hides changes in every repository that has them, under one name. New (untracked) files are stashed too, ignored files (for example `node_modules`) are left alone.

```bash
amr stash -m login-fix          # stash under the name login-fix
amr stash . -m login-fix        # only in the repository you are in
amr stash list                  # amr stashes by name and the repositories that have them
amr stash pop login-fix         # restore login-fix everywhere and remove it
amr stash pop                   # restore the latest amr stash in each repository
amr stash drop login-fix        # remove login-fix without restoring
```

This is handy when one task spans several repositories: stash all of them with one name and get them back later with one command.

Without `-m`, the name is generated from the current time, for example `amr-2026-10-15-143005`, and is printed after stashing.

amr finds stashes by name, so you don't need numbers like `stash@{2}`. Your regular `git stash` entries are not affected.

## Working in parallel

`clone`, `status`, `pull`, `sync` and `run` process up to 8 repositories at the same time. Change it with the `--jobs` (`-j`) flag:

```bash
amr pull -j 16          # faster when there are many repositories
amr run -j 1 -- make    # strictly one by one
```

While `pull`, `clone` and `sync` are running, a progress line is shown in the terminal. When the output goes to a file or another command (`amr pull | grep ...`), there is no progress and no colors, only the plain result.

## Workspace config

The config is stored in `.amr/amr.yaml`. `amr add` fills it in, taking the URL from each repository's `origin`:

```yaml
version: 1
repositories:
  - path: frontend
    url: git@github.com:your-org/frontend.git
  - path: services/api
    url: git@github.com:your-org/api.git
    groups: [backend]
```

| Field | Meaning |
| --- | --- |
| `path` | path from the workspace root, can be nested |
| `url` | URL to clone from |
| `groups` | groups, optional |

A repository without `origin` is saved without `url`, and `amr clone` skips it.

You can edit the config by hand, for example to change a path or groups.

## Workspace for a team

It is convenient to keep the config in its own Git repository, for example `your-org/acme-workspace`. The name can be anything.

`amr init` creates a `.gitignore` by itself, so that only the config goes into this repository and the cloned repositories stay out of it. If the folder already has a `.gitignore` or is a Git repository itself, the file is not created:

```gitignore
/*
!/.amr/
!/.gitignore
!/README.md
```

A new teammate gets the whole layout with three commands:

```bash
git clone git@github.com:your-org/acme-workspace.git acme
cd acme
amr sync
```

After that, just run `amr sync` from time to time. It updates the config first, so repositories added by teammates are cloned too.

## Development

Requires Go and [Task](https://taskfile.dev).

```bash
task run -- --help   # run from source
task build           # build the binary into build/amr
task check           # formatting, dependencies, build, tests
task test            # run tests
task bench           # run benchmarks
```
