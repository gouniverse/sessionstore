# Session Store

Stores session to a database table.

## Installation
```
go get -u github.com/gouniverse/sessionstore
```

## Setup

```
sessionStore = sessionstore.NewStore(sessionstore.WithGormDb(databaseInstance), sessionstore.WithTableName("my_session"), sessionstore.WithAutoMigrate(true))

go sessionStore.ExpireSessionGoroutine()
```
