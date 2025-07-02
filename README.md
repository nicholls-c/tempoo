# tempoo-go

- [tempoo-go](#tempoo-go)
  - [Links](#links)
  - [Auth](#auth)
  - [Use](#use)
    - [Using a binary](#using-a-binary)
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

## Use

### Using a binary

```sh
chmod +x ./dist/tempoo
sudo mv ./tempoo /usr/local/bin
```

### Add worklog

```sh
# uncompiled
go run cmd/main.go add-worklog --issue-key INF-88 --time 3h
go run cmd/main.go add-worklog -i INF-88 -t 1h
```

<br>

### Remove worklogs:

```sh
go run cmd/main.go remove-worklog --issue-key INF-88
```

<br>

### Flags

Supply `-d` to any command to get verbose output.

<br>

## Build

```sh
# validate
goreleaser check

# build
goreleaser release --snapshot --clean
```

<br>

## Test

```sh
go test -v ./cmd/...

go test -bench=. ./cmd/...

go test -cover ./cmd/...
```

<br>