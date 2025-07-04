# tempoo

- [tempoo](#tempoo)
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

<br>

## Get

Navigate to the [release](https://github.com/nicholls-c/tempoo/releases) page for a list of all avaialable releases. Use the latest version by default unless told otherwise.

### Linux/WSL

1. Download release using `gh` cli:
   ```sh
   gh release download v0.1.4 --pattern "tempoo" -R nicholls-c/tempoo-go --clobber
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
   gh release download v0.1.2 --pattern "tempoo.exe" -R nicholls-c/tempoo-go --clobber
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
export JIRA_EMAIL=<firstname.lastname>@commify.com
export JIRA_API_TOKEN=myapitoken
```

#### Windows

```powershell
[string]$env:JIRA_EMAIL = "<firstname.lastname>@commify.com"
[string]$env:$JIRA_API_TOKEN = "myapitoken"
```

<br>

### Add worklog

```sh
# defaults to today
tempoo add-worklog --issue-key INF-88 --time 3h

# for specified date
tempoo add-worklog --issue-key INF-88 --time 8h --date 01.07.2025
```

<br>

### Remove worklogs

```sh
tempoo remove-worklog --issue-key INF-88 --verbose
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