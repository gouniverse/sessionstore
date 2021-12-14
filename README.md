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

## Usage

```
sessionKey  := "ABCDEFG"
sessionExpireSeconds = 2*60*60

// Create new / update existing session
sessionStore.Set(sessionKey, sessionValue, sessionExpireSeconds)

// Get session value, or default if not found
sessionValue := sessionStore.Get(sessionKey, defaultValue)
```

## Changelog

2021.12.14 - Added support for DB dialects

2021.12.14 - Removed GORM dependency and moved to the standard library