# s1gu framework

[![Build Status](https://travis-ci.com/jiharal/s1gu.svg?branch=master)](https://travis-ci.com/jiharal/s1gu)
[![Coverage Status](https://coveralls.io/repos/github/jiharal/s1gu/badge.svg?branch=master)](https://coveralls.io/github/jiharal/s1gu?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jiharal/s1gu)](https://goreportcard.com/report/github.com/jiharal/s1gu)

s1gu is a RESTful And GraphQL framework for the rapid development of Go applications including APIs.

## s1gu commands

```cmd
  A GraphQL and RESTful API Framework Go

  Usage:s
    s1gu [command]

  Available Commands:
    help        Help about any command
    model       Create model application
    new         Create new project
    router      Create router application

  Flags:
    -h, --help   help for s1gu

  Use "s1gu [command] --help" for more information about a command.
```

## Create new project

### Command

```cmd
$ s1gu new MyApp
```

### Response

```cmd
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/cmd/cmd.go
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/model/model.user.go
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/router/init.go
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/router/handler.go
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/router/graphql.user.go
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/router/restful.user.go
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/main.go
  create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/.myapp.toml
New application successfully created!
```

## Create new model

```cmd
Usage:
  s1gu model [model name]
```

- `[model name]` use the table name in your database

### Command

```cmd
$ cd MyApp
$ s1gu model access
```

### Response

```cmd
2018/09/20 01:00:23 Do you want to add it? [Yes|No]
yes
create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/model/model.access.go
```

## Create new router

```cmd
Usage:
  s1gu router [router name] [graphql or rest]
```

- `[router name]` use the table name in your database
- `[graphql or rest]` use one of the commands between `graphql` and `rest`

### Command

```cmd
$ s1gu router access rest
```

### Response RESTful API Base

```cmd
2018/09/20 01:10:00 Do you want to add it? [Yes|No]
y
create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/router/restful.access.go
```

### Command

```cmd
$ s1gu router access graphql
```

### Response GraphQL API Base

```cmd
2018/09/20 01:13:34 Do you want to add it? [Yes|No]
y
create	 /Users/zzz/go/src/github.com/jiharal/s1gu/MyApp/router/graphql.access.go
```

## Adding handler in `cmd/cmd.go`

```go
...
// Access API
r.HandleFunc("/access", router.GetAllAccess).Methods("GET")
r.HandleFunc("/access/{id}", router.GetOneAccess).Methods("GET")
r.HandleFunc("/access", router.InsertAccess).Methods("POST")
r.HandleFunc("/access/{id}", router.UpdateAccess).Methods("POST")
r.HandleFunc("/access/{id}", router.DeleteAccess).Methods("DELETE")
...
```

## Config your APP

`.myapp.toml`

```toml

[app]
name = "MyApp"
version = "0.0.1"
port = 9100

[database]
host = "localhost"
port = 26257
username = "root"
password = ""
name = "myapp_db"
sslmode = "disable"

[cache]
host = "localhost"
port = 6379
password = ""
max_idle = 100
idle_timeout = 5
enabled = true
expire_time = 60
idempotency_expiry = 86400
```
