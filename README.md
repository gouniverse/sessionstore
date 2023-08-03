# Session Store <a href="https://gitpod.io/#https://github.com/gouniverse/sessionstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

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
	TimeoutSeconds:     3600, // 1 hour
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})

go sessionStore.ExpireSessionGoroutine()
```

## Methods

- AutoMigrate() error - automigrate (creates) the session table
- DriverName(db *sql.DB) string - finds the driver name from database
- EnableDebug(debug bool) - enables / disables the debug option
- ExpireSessionGoroutine() error - deletes the expired session keys
- Delete(sessionKey string) (bool, error)  - Delete deletes a session
- FindByKey(sessionKey string) (*Session, error) - FindByKey finds a session by key
- Get(sessionKey string, valueDefault string) (string, error) - Gets the session value as a string
- GetAny(key string, valueDefault any) (any, error) - attempts to parse the value as interface, use with SetAny
- GetMap(key string, valueDefault map[string]any) (map[string]any, error) - attempts to parse the value as map[string]any, use with SetMap
- Has(sessionKey string) (bool, error) - Checks if a session by key exists
- Set(sessionKey string, value string, seconds int64) error - Set sets a key in store
- SetAny(key string, value any, seconds int64) error - convenience method which saves the supplied interface value, use GetAny to extract
- MergeMap(key string, mergeMap map[string]any, seconds int64) error - updates an existing map
- SetMap(key string, value map[string]any, seconds int64) error - convenience method which saves the supplied map, use GetMap to extract

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
// Store interface value
sessionStore.SetAny(sessionKey, sessionValue, sessionExpireSeconds)

// Get interface value
value := sessionStore.GetAny(sessionKey, defaultValue)



// Example
value := map[string]string{
  "key1": "value1",
  "key2": "value2",
  "key3": "value3",
}
isSaved, err := store.SetJSON("mykey", value, 5*60)

if !isSaved {
  log.Fatal("Set failed: " + err.Error())
}

result, err := store.GetJSON("mykey", "{}")

if err != nil {
  log.Fatal("Get failed: " + err.Error())
}

var res = map[string]string{}
for k, v := range result.(map[string]interface{}) {
  res[k] = v.(string)
}

log.Println(res["key1"])
```



## Changelog

2023.08.03 - Renamed "SetJSON", "GetJSON" methods to "SetAny", "GetAny"

2023.08.03 - Added "SetMap", "GetMap", "MergeMap" methods

2022.12.06 - Changed store setup to use struct

2022.01.01 - Added "Has" method

2021.12.15 - Added LICENSE

2021.12.15 - Added test badge

2021.12.15 - Added SetJSON GetJSON

2021.12.14 - Added support for DB dialects

2021.12.14 - Removed GORM dependency and moved to the standard library
