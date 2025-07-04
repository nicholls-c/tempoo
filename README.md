![Tempoo](./docs/images/tempoo.png)

[![ci](https://github.com/nicholls-c/tempoo/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/nicholls-c/tempoo/actions/workflows/ci.yml)

[![Latest Release](https://img.shields.io/github/v/release/nicholls-c/tempoo)](https://github.com/nicholls-c/tempoo/releases/latest)

Automate the awful things.

---

<br>

- [Links](#links)
- [Get](#get)
  - [Linux/WSL](#linuxwsl)
  - [Windows](#windows)
- [Use](#use)
  - [Authenticate](#authenticate)
    - [Linux/WSL](#linuxwsl-1)
    - [Windows](#windows-1)
  - [Add worklog](#add-worklog)
  - [Remove worklogs](#remove-worklogs)
  - [List worklogs](#list-worklogs)
  - [Show app version](#show-app-version)
  - [Debug](#debug)
- [Contributing](#contributing)
  - [Pre Commit](#pre-commit)
  - [Build](#build)
  - [Test](#test)
  - [Creating Release](#creating-release)


## Links

- [Kong](https://github.com/alecthomas/kong)
- [JIRA REST API](https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/)

<br>

## Get

Navigate to the [release](https://github.com/nicholls-c/tempoo/releases) page for a list of all avaialable releases. Use the latest version by default unless told otherwise.

### Linux/WSL

1. Download release using `gh` cli:
   ```sh
   gh release download v0.1.4 --pattern "tempoo" -R nicholls-c/tempoo --clobber
   ```
2. Make exectable:
   ```sh
   sudo chmod +x ./tempoo
   ```
3. Move to PATH:
   ```sh
   sudo mv ./tempoo /usr/local/bin
   ```

### Windows

1. Download release using `gh` cli:
   ```sh
   gh release download v0.1.2 --pattern "tempoo.exe" -R nicholls-c/tempoo --clobber
   ```
2. Validate:
   ```pwsh
   ./tempoo.exe -h
   ```

Or just download it from the releases pages manually :blush:

<br>

## Use

### Authenticate

Expose user email address and jira api token as env vars.

#### Linux/WSL

```sh
export JIRA_EMAIL=<firstname.lastname>@emaildomain.com
export JIRA_API_TOKEN=myapitoken
```

#### Windows

```powershell
[string]$env:JIRA_EMAIL = "<firstname.lastname>@emaildomain.com"
[string]$env:$JIRA_API_TOKEN = "myapitoken"
```

<br>

### Add worklog

```sh
# defaults to today
tempoo add-worklog --issue-key INF-88 --time 3

# for specified date
tempoo add-worklog --issue-key INF-88 --time 8 --date 01.07.2025 --verbose
```

<br>

### Remove worklogs

```sh
tempoo remove-worklogs --issue-key INF-88 --verbose
```

<br>

### List worklogs

```sh
tempoo list-worklogs -i INF-88
```

<br>

### Show app version

```sh
tempoo version
```

<br>

### Debug

Supply `--verbose` to any command to get verbose debug output.

<br>

## Contributing

### Pre Commit

A pre-commit config is in place.

Install pre-commit to use.

```sh
pipx install pre-commit
pre-commit install
```

<br>

### Build

```sh
# validate
goreleaser check

# build
goreleaser release --snapshot --clean

sudo mv dist/tempoo-go_linux_amd64_v1/tempoo /usr/local/bin/
tempoo version # should have SNAPSHOT in the version string
```

<br>

### Test

> **Warning**
> All tests were cursor'd

```sh
go test -v ./...

go test -bench=. ./...

go test -cover ./...
```

<br>

### Creating Release

Release are built from tags, see [workflow](./.github/workflows/ci.yml).

```sh
git tag v1.0.0
git push origin v1.0.0
```
