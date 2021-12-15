# Session Store

[![Tests Status](https://github.com/gouniverse/cachestore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/sessionstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/cachestore)](https://goreportcard.com/report/github.com/gouniverse/sessionstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/cachestore)](https://pkg.go.dev/github.com/gouniverse/sessionstore)

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
value := sessionStore.Get(sessionKey, defaultValue)

// Delete session
isDeleted, err := sessionStore.Delete(sessionKey)
```



```
// Store JSON value
sessionStore.SetJSON(sessionKey, sessionValue, sessionExpireSeconds)

// Get JSON value
value := sessionStore.GetJSON(sessionKey, defaultValue)
```



## Changelog

2021.12.15 - Added SetJSON GetJSON

2021.12.14 - Added support for DB dialects

2021.12.14 - Removed GORM dependency and moved to the standard library
