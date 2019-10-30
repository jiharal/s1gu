# AuthscureGo framework

authscure-go is a RESTful And GraphQL framework for the rapid development of Go applications including APIs.

## authscure-go commands

  ```cmd
    A GraphQL and RESTful API Framework Go

    Usage:
      authscure-go [command]

    Available Commands:
      help        Help about any command
      model       Create model application
      new         Create new project
      router      Create router application

    Flags:
      -h, --help   help for authscure-go

    Use "authscure-go [command] --help" for more information about a command.
  ```

## Create new project

### Command

```cmd
$ authscure-go new MyApp
```

### Response

  ```cmd
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/cmd/cmd.go
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/model/model.user.go
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/router/init.go
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/router/handler.go
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/router/graphql.user.go
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/router/restful.user.go
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/main.go
    create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/.myapp.toml
  New application successfully created!
  ```

## Create new model

  ```cmd
  Usage:
    authscure-go model [model name]
  ```
  - `[model name]` use the table name in your database

### Command

  ```cmd
  $ cd MyApp
  $ authscure-go model access
  ```

### Response

  ```cmd
  2018/09/20 01:00:23 Do you want to add it? [Yes|No]
  yes
	create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/model/model.access.go
  ```

## Create new router

  ```cmd
  Usage:
    authscure-go router [router name] [graphql or rest]
  ```
  - `[router name]` use the table name in your database
  - `[graphql or rest]` use one of the commands between `graphql` and `rest`

### Command

  ```cmd
  $ authscure-go router access rest
  ```

### Response RESTful API Base

  ```cmd
  2018/09/20 01:10:00 Do you want to add it? [Yes|No]
  y
	create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/router/restful.access.go
  ```

### Command

  ```cmd
  $ authscure-go router access graphql
  ```

### Response GraphQL API Base

  ```cmd
  2018/09/20 01:13:34 Do you want to add it? [Yes|No]
  y
	create	 /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/router/graphql.access.go
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

## Testing
  
### CockroachDB

  ```cmd
  $ cockroach start --insecure
  ```

#### Create database

  ```cmd
  $ cockroach sql --insecure
   > create database myapp_db;
   > use myapp_db;
  ```

#### Create table `user` and `access`

  ```sql
  > CREATE TABLE "user" (
    id             UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    name      STRING  NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by     UUID    NOT NULL,
    updated_at     TIMESTAMPTZ,
    updated_by     UUID
  );

  > CREATE TABLE access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID NOT NULL,
    updated_at TIMESTAMPTZ,
    updated_by UUID
  );
  ```
#### Insert data to table access
  
  ```sql
  > INSERT INTO "user"
    (id, name, created_by) 
    VALUES ('c01c2cbb-e6c9-4d45-a486-2929d5404ea1', 'jihar', 'c01c2cbb-e6c9-4d45-a486-2929d5404eb2');
  
  > INSERT INTO access(id, name, created_by) VALUES
    ('02842d9a-979d-4cad-b2eb-0dd131c11e91', 'CATEGORY_VIEW', '819a6572-d825-4dc4-8d0a-71177e62e795'),
    ('9a2c64de-79c1-4430-8ed1-78797835e761', 'CATEGORY_CREATE', '819a6572-d825-4dc4-8d0a-71177e62e795'),
    ('4cbee14b-89b8-4024-8e28-00e9e0c5ceea', 'CATEGORY_UPDATE', '819a6572-d825-4dc4-8d0a-71177e62e795');
  ```

### Redis

  ```cmd
  redis-server
  ```
### App

  ```cmd
  $ go run main.go
  using config file:  /Users/zzz/go/src/github.com/AuthScureDevelopment/authscure-go/MyApp/.myapp.toml
  Listening on http://localhost:9100
  ```

### Browser

  ```url
    http://localhost:9100/access
  ```
  ![img](https://user-images.githubusercontent.com/21150538/45773062-8dbc7800-bc73-11e8-8e6e-9c9e2718e072.png)