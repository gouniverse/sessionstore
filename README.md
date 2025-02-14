# Session Store
[![Tests Status](https://github.com/gouniverse/sessionstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/sessionstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/sessionstore)](https://goreportcard.com/report/github.com/gouniverse/sessionstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/sessionstore)](https://pkg.go.dev/github.com/gouniverse/sessionstore)

Stores session to a database table.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## 🌏  Open in the Cloud 
Click any of the buttons below to start a new development environment to demo or contribute to the codebase without having to install anything on your machine:

[![Open in VS Code](https://img.shields.io/badge/Open%20in-VS%20Code-blue?logo=visualstudiocode)](https://vscode.dev/github/gouniverse/sessionstore)
[![Open in Glitch](https://img.shields.io/badge/Open%20in-Glitch-blue?logo=glitch)](https://glitch.com/edit/#!/import/github/gouniverse/sessionstore)
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/gouniverse/sessionstore)
[![Edit in Codesandbox](https://codesandbox.io/static/img/play-codesandbox.svg)](https://codesandbox.io/s/github/gouniverse/sessionstore)
[![Open in StackBlitz](https://developer.stackblitz.com/img/open_in_stackblitz.svg)](https://stackblitz.com/github/gouniverse/sessionstore)
[![Open in Repl.it](https://replit.com/badge/github/withastro/astro)](https://replit.com/github/gouniverse/sessionstore)
[![Open in Codeanywhere](https://codeanywhere.com/img/open-in-codeanywhere-btn.svg)](https://app.codeanywhere.com/#https://github.com/gouniverse/sessionstore)
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/gouniverse/sessionstore)



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

go sessionStore.SessionExpiryGoroutine()
```

## Methods

- AutoMigrate() error - automigrate (creates) the session table
- DriverName(db *sql.DB) string - finds the driver name from database
- EnableDebug(debug bool) - enables / disables the debug option
- SessionExpiryGoroutine() error - deletes the expired session keys

## Usage

```go
sessionKey  := "ABCDEFG"
sessionExpireSeconds = 2*60*60

session := NewSession().
  SetKey(sessionKey).
  SetValue(sessionValue).
  SetUserID(userID).
  SetUserAgent(r.UserAgent()).
  SetIPAddress(r.RemoteAddr).
  SetExpiresAt(carbon.Now(carbon.UTC).AddSeconds(sessionExpireSeconds).ToDateTimeString(carbon.UTC))

// Create new
err := sessionStore.SessionCreate(session)

// Get session value, or default if not found
session, err := sessionStore.SessionFindByKey(sessionKey)

// Update session
session.SetValue(newSessionValue)
session.SetExpiresAt(carbon.Now(carbon.UTC).AddMinutes(60).ToDateTimeString(carbon.UTC))
err := sessionStore.SessionUpdate(session)

// Delete session
err := sessionStore.SessionDeleteByKey(sessionKey)
```


## Changelog

2025.01.05 - Added "SessionExtend" method

2024.12.11 - Removed old API, extended interface

2024.09.08 - Added options (UserID, UserAgent, IPAddress)

2024.01.03 - Added "Extend" method

2023.08.03 - Renamed "SetJSON", "GetJSON" methods to "SetAny", "GetAny"

2023.08.03 - Added "SetMap", "GetMap", "MergeMap" methods

2022.12.06 - Changed store setup to use struct

2022.01.01 - Added "Has" method

2021.12.15 - Added LICENSE

2021.12.15 - Added test badge

2021.12.15 - Added SetJSON GetJSON

2021.12.14 - Added support for DB dialects

2021.12.14 - Removed GORM dependency and moved to the standard library
