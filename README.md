# txtract

[![License](https://img.shields.io/github/license/FollowTheProcess/txtract)](https://github.com/FollowTheProcess/txtract)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/txtract)](https://goreportcard.com/report/github.com/FollowTheProcess/txtract)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/txtract?logo=github&sort=semver)](https://github.com/FollowTheProcess/txtract)
[![CI](https://github.com/FollowTheProcess/txtract/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/txtract/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/FollowTheProcess/txtract/branch/main/graph/badge.svg)](https://codecov.io/gh/FollowTheProcess/txtract)

A CLI to interact with txtar archive files 📂

## Project Description

[txtar] is an incredibly useful tiny archive format, able to trivially store a miniature filesystem in a single plain text file. It particularly shines for storing test cases!

Most of the time a txtar archive is created manually, but what if you already have a directory full of stuff and you want to instead store them in txtar. Or you have a txtar
archive that you want to instantly replicate on your filesystem!

That's where `txtract` comes in 🚀

![quickstart](https://github.com/FollowTheProcess/txtract/raw/main/docs/img/demo.gif)

## Installation

Compiled binaries for all supported platforms can be found in the [GitHub release]. There is also a [homebrew] tap:

```shell
brew install -cask FollowTheProcess/tap/txtract
```

## Quickstart

Recursively zip up the contents of a directory into a single txtar file named yourdirectory.txtar

```shell
txtract zip ./yourdirectory
```

Or go the other way, unzip a txtar file back into your filesystem:

```shell
txtract unzip yourdirectory.txtar
```

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[GitHub release]: https://github.com/FollowTheProcess/txtract/releases
[homebrew]: https://brew.sh
[txtar]: https://pkg.go.dev/golang.org/x/tools/txtar
