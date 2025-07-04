# tempoo-go

- [tempoo-go](#tempoo-go)
  - [Links](#links)
  - [Auth](#auth)
  - [Use](#use)
    - [Using a binary](#using-a-binary)
      - [Linux](#linux)
      - [Windows](#windows)
    - [Add worklog](#add-worklog)
    - [Remove worklogs:](#remove-worklogs)
    - [Flags](#flags)
  - [Build](#build)
  - [Test](#test)


## Links

- [Kong](https://github.com/alecthomas/kong)
- [Jira Example](https://esendex.atlassian.net/browse/INF-88)

<br>

## Auth

Expose user email address and Jira API token:

```sh
export JIRA_EMAIL=christian.nicholls@commify.com
export JIRA_API_TOKEN=asdasdasdasdasdasdasdasdasd
```

<br>

## Get

### Linux/WSL

1. Download release using `gh` cli:
   ```sh
   gh release download v0.1.2 --pattern "tempoo" -R nicholls-c/tempoo-go --clobber
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
   gh release download v0.1.2 --pattern "tempoo" -R nicholls-c/tempoo-go --clobber
   ```
2. Validate:
   ```pwsh
   ./tempoo.exe -h
   ```

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

> [!CAUTION] All tests were cursor'd.

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