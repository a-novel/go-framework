# Framework

Common libraries for the backend services.

```bash
go get github.com/a-novel/go-framework
```

## Installation

You must run this command at least once, for the entire project.

```bash
docker compose up -d
```

## Commands

Connect to the database:

```bash
make db
# psql (14.6 (Homebrew), server 14.7 (Debian 14.7-1.pgdg110+1))
# Type "help" for help.
# 
# agora=#
```

Connect to the test database:

```bash
make db-test
# psql (14.6 (Homebrew), server 14.7 (Debian 14.7-1.pgdg110+1))
# Type "help" for help.
# 
# agora_test=>
```

Run tests

```bash
make test
```
