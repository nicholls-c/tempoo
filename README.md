# tempoo-go

- [tempoo-go](#tempoo-go)
  - [Auth:](#auth)
  - [Use](#use)
    - [Add worklog](#add-worklog)
    - [Remove worklogs:](#remove-worklogs)


## Auth:

Expose user email address and Jira API token:

```sh
export JIRA_EMAIL=christian.nicholls@commify.com
export JIRA_API_TOKEN=asdasdasdasdasdasdasdasdasd
```

<br>

## Use

### Add worklog

```sh
go run cmd/main.go add-worklog --issue-key INF-88 --time 3h
go run cmd/main.go add-worklog -i INF-88 -t 1h
```

<br>

### Remove worklogs:

```sh

```

<br>