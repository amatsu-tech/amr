# Changelog

## [0.3.0](https://github.com/amatsu-tech/amr/compare/v0.2.0...v0.3.0) (2026-10-13)


### Features

* choose repositories by path and group ([7c861be](https://github.com/amatsu-tech/amr/commit/7c861befb0befda1b7c56879673947b24f3bd29d))
* create a workspace .gitignore on amr init ([ebba1a8](https://github.com/amatsu-tech/amr/commit/ebba1a851a4ced01f1011e017f10316507dd99fc))
* sync the workspace config and all repositories with amr sync ([a736d81](https://github.com/amatsu-tech/amr/commit/a736d81180a8e76ddeda9f5f4dffbd3cff8343f7))


### Refactoring

* share group flag text and empty workspace message ([d2bfd31](https://github.com/amatsu-tech/amr/commit/d2bfd31246c289e236e48c3b4085329a7a280cd9))

## [0.2.0](https://github.com/amatsu-tech/amr/compare/v0.1.0...v0.2.0) (2026-10-13)


### Features

* add every found repository with amr add --all ([476c4d9](https://github.com/amatsu-tech/amr/commit/476c4d9a843752d7fc066aef51a2388c05691bd5))
* process repositories in parallel with --jobs ([b86e6d8](https://github.com/amatsu-tech/amr/commit/b86e6d8dc825ca79edaf32210a3d5891983dd194))
* show progress while pulling and cloning ([de50555](https://github.com/amatsu-tech/amr/commit/de505555853f9b91e71e2710598915da995b3768))
* suggest amr clone in status for missing repositories ([1901e2c](https://github.com/amatsu-tech/amr/commit/1901e2cdbac926e5f3ae690ea11409b8fa7aa082))


### Bug Fixes

* print output to stdout and drop colors when not a terminal ([7b3f5da](https://github.com/amatsu-tech/amr/commit/7b3f5da0d9014e808b49536b451334c04b637b6d))

## 0.1.0 (2026-10-11)


### Features

* add output package ([e73672d](https://github.com/amatsu-tech/amr/commit/e73672d4edeaefde64a07f69a44b18ad9268bf38))
* add workspace management commands ([9057b31](https://github.com/amatsu-tech/amr/commit/9057b319ee917b4a16eda78bb814a7d38088e179))
* show search path when workspace is not found ([5067fb7](https://github.com/amatsu-tech/amr/commit/5067fb7fcb4680dbe99e1059339bd7352a56735a))
* support nested repository paths and add clone command ([7f52258](https://github.com/amatsu-tech/amr/commit/7f52258caf09f1d18e45723ad028982398c2e226))
