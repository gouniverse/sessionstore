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

```go
sessionStore = sessionstore.NewStore(sessionstore.NewStoreOptions{
	DB:                 databaseInstance,
	SessionTableName:   "my_session",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})

go sessionStore.ExpireSessionGoroutine()
```

## Usage

```go
sessionKey  := "ABCDEFG"
sessionExpireSeconds = 2*60*60

// Create new / update existing session
sessionStore.Set(sessionKey, sessionValue, sessionExpireSeconds)

// Get session value, or default if not found
value := sessionStore.Get(sessionKey, defaultValue)

// Delete session
isDeleted, err := sessionStore.Delete(sessionKey)
```



```go
// Store JSON value
sessionStore.SetJSON(sessionKey, sessionValue, sessionExpireSeconds)

// Get JSON value
value := sessionStore.GetJSON(sessionKey, defaultValue)



// Example
value := map[string]string{
  "key1": "value1",
  "key2": "value2",
  "key3": "value3",
}
isSaved, err := store.SetJSON("mykey", value, 5*60)

if !isSaved {
  log.Fatal("Set JSON failed: " + err.Error())
}

result, err := store.GetJSON("mykey", "{}")

if err != nil {
  log.Fatal("Get JSON failed: " + err.Error())
}

var res = map[string]string{}
for k, v := range result.(map[string]interface{}) {
  res[k] = v.(string)
}

log.Printls(res["key1"])
```



## Changelog

2022.12.06 - Changed store setup to use struct

2022.01.01 - Added "Has" method

2021.12.15 - Added LICENSE

2021.12.15 - Added test badge

2021.12.15 - Added SetJSON GetJSON

2021.12.14 - Added support for DB dialects

2021.12.14 - Removed GORM dependency and moved to the standard library
