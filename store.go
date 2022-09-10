package sessionstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"     // importing mysql dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"  // importing postgres dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"   // importing sqlite3 dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver" // importing sqlserver dialect
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/uid"
)

// Store defines a session store
type Store struct {
	sessionTableName   string
	db                 *sql.DB
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debug              bool
}

// StoreOption options for the session store
type StoreOption func(*Store)

// WithAutoMigrate sets the table name for the cache store
func WithAutoMigrate(automigrateEnabled bool) StoreOption {
	return func(s *Store) {
		s.automigrateEnabled = automigrateEnabled
	}
}

// WithDb sets the database for the setting store
func WithDb(db *sql.DB) StoreOption {
	return func(s *Store) {
		s.db = db
		s.dbDriverName = s.DriverName(s.db)
	}
}

// WithDebug prints the SQL queries
func WithDebug(debug bool) StoreOption {
	return func(s *Store) {
		s.debug = debug
	}
}

// WithTableName sets the table name for the session store
func WithTableName(sessionTableName string) StoreOption {
	return func(s *Store) {
		s.sessionTableName = sessionTableName
	}
}

// NewStore creates a new session store
func NewStore(opts ...StoreOption) (*Store, error) {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.sessionTableName == "" {
		return nil, errors.New("session store: sessionTableName is required")
	}

	if store.debug {
		log.Println(store.dbDriverName)
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() error {
	sql := st.SQLCreateTable()

	_, err := st.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DriverName finds the driver name from database
func (st *Store) DriverName(db *sql.DB) string {
	dv := reflect.ValueOf(db.Driver())
	driverFullName := dv.Type().String()
	if strings.Contains(driverFullName, "mysql") {
		return "mysql"
	}
	if strings.Contains(driverFullName, "postgres") || strings.Contains(driverFullName, "pq") {
		return "postgres"
	}
	if strings.Contains(driverFullName, "sqlite") {
		return "sqlite"
	}
	if strings.Contains(driverFullName, "mssql") {
		return "mssql"
	}
	return driverFullName
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debug = debug
}

// ExpireSessionGoroutine - soft deletes expired cache
func (st *Store) ExpireSessionGoroutine() error {
	i := 0
	for {
		i++
		if st.debug {
			log.Println("Cleaning expired sessions...")
		}
		sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.sessionTableName).Where(goqu.C("expires_at").Lt(time.Now())).Delete().ToSQL()

		if st.debug {
			log.Println(sqlStr)
		}

		_, err := st.db.Exec(sqlStr)
		if err != nil {
			if err == sql.ErrNoRows {
				// Looks like this is now outdated for sqlscan
				return nil
			}
			if sqlscan.NotFound(err) {
				return nil
			}
			log.Println("Cache Store. ExpireSessionGoroutine. Error: ", err)
			return nil
		}

		time.Sleep(60 * time.Second) // Every minute
	}
}

// Delete deletes a session
func (st *Store) Delete(sessionKey string) (bool, error) {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.sessionTableName).Where(goqu.C("session_key").Eq(sessionKey)).Delete().ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			// Looks like this is now outdated for sqlscan
			return false, nil
		}

		if sqlscan.NotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// FindByKey finds a session by key
func (st *Store) FindByKey(sessionKey string) *Session {
	// key exists, expires is < now, deleted null
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.sessionTableName).Where(goqu.C("session_key").Eq(sessionKey), goqu.C("expires_at").Gt(time.Now()), goqu.C("deleted_at").IsNull()).Select("*").ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var session Session
	err := sqlscan.Get(context.Background(), st.db, &session, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// Looks like this is now outdated for sqlscan
			return nil
		}
		if sqlscan.NotFound(err) {
			return nil
		}
		log.Println("CacheStore. FindBykey. Error: ", err)
		return nil
	}

	return &session
}

// Gets the session value as a string
func (st *Store) Get(sessionKey string, valueDefault string) string {
	session := st.FindByKey(sessionKey)

	if session != nil {
		return session.Value
	}

	return valueDefault
}

// GetJSON attemots to parse the value as JSON, use with SetJSON
func (st *Store) GetJSON(key string, valueDefault interface{}) (interface{}, error) {
	session := st.FindByKey(key)

	if session != nil {
		jsonValue := session.Value
		var e interface{}
		jsonError := json.Unmarshal([]byte(jsonValue), &e)
		if jsonError != nil {
			return valueDefault, jsonError
		}

		return e, nil
	}

	return valueDefault, nil
}

// Has finds if a session by key exists
func (st *Store) Has(sessionKey string) (bool, error) {
	// key exists, expires is < now, deleted null
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.sessionTableName).Where(goqu.C("session_key").Eq(sessionKey), goqu.C("expires_at").Gt(time.Now()), goqu.C("deleted_at").IsNull()).Select(goqu.COUNT("*")).As("count").ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var count int
	err := sqlscan.Get(context.Background(), st.db, &count, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// Looks like this is now outdated for sqlscan
			return false, nil
		}
		if sqlscan.NotFound(err) {
			return false, nil
		}
		log.Println("CacheStore. Error: ", err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// Set sets a key in store
func (st *Store) Set(sessionKey string, value string, seconds int64) (bool, error) {
	session := st.FindByKey(sessionKey)

	expiresAt := time.Now().Add(time.Second * time.Duration(seconds))

	var sqlStr string
	if session == nil {
		var newSession = Session{
			ID:        uid.MicroUid(),
			Key:       sessionKey,
			Value:     value,
			ExpiresAt: &expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.sessionTableName).Rows(newSession).ToSQL()
	} else {
		fields := map[string]interface{}{}
		fields["session_value"] = value
		fields["expires_at"] = &expiresAt
		fields["updated_at"] = time.Now()
		sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Update(st.sessionTableName).Set(fields).ToSQL()
	}

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return false, err
	}

	return true, nil
}

// SetJSON convenience method which saves the supplied value as JSON, use GetJSON to extract
func (st *Store) SetJSON(key string, value interface{}, seconds int64) (bool, error) {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return false, jsonError
	}

	return st.Set(key, string(jsonValue), seconds)
}

// // SessionGetKey gets a key from sessiom
// func SessionGetKey(sessionKey string, key string, valueDefault string) string {
// 	session := SessionFindByToken(sessionKey)

// 	if session == nil {
// 		return valueDefault
// 	}

// 	kv := hashmap.New()
// 	err := kv.FromJSON([]byte(session.Value))
// 	if err != nil {
// 		return valueDefault
// 	}

// 	value, _ := kv.Get(key)
// 	if value != nil {
// 		return fmt.Sprintf("%v", value)
// 	}

// 	return valueDefault
// }

// Empty removes all keys from the sessiom
// func (st *Store) Empty(sessionKey string) (bool, error) {
// 	session := st.FindByKey(sessionKey)

// 	kv := hashmap.New()

// 	if session == nil {
// 		return true, nil
// 	}

// 	json, err := kv.ToJSON()

// 	if err != nil {
// 		return false, err
// 	}

// 	session.Value = string(json)

// 	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).Update(st.sessionTableName).Set(session).ToSQL()

// 	if st.debug {
// 		log.Println(sqlStr)
// 	}

// 	_, errExec := st.db.Exec(sqlStr)

// 	if errExec != nil {
// 		return false, errExec
// 	}

// 	return true, nil
// }

// SetKey sets a single key into sessiom
// func (st *Store) SetKey(sessionKey string, key string, value string) (bool, error) {
// 	session := st.FindBySessionKey(sessionKey)

// 	kv := hashmap.New()

// 	if session == nil {
// 		isOk, err := st.Set(sessionKey, "{}", 2000)
// 		if isOk == false {
// 			return false, err
// 		}
// 		session = st.FindBySessionKey(sessionKey)
// 	}

// 	log.Println(value)

// 	err := kv.FromJSON([]byte(session.Value))
// 	if err != nil {
// 		return false, err
// 	}

// 	kv.Put(key, value)
// 	json, err := kv.ToJSON()

// 	if err != nil {
// 		return false, err
// 	}

// 	session.Value = string(json)

// 	seconds := session.ExpiresAt.Unix() - time.Now().Unix()
// 	return st.Set(sessionKey, string(json), seconds)
// }

// RemoveKey removes a key from sessiom
// func (st *Store) RemoveKey(sessionKey string, key string) (bool, error) {
// 	session := st.FindBySessionKey(sessionKey)

// 	kv := hashmap.New()

// 	if session == nil {
// 		return true, nil
// 	}

// 	err := kv.FromJSON([]byte(session.Value))
// 	if err != nil {
// 		return false, err
// 	}

// 	kv.Remove(key)

// 	json, err := kv.ToJSON()

// 	if err != nil {
// 		return false, err
// 	}

// 	log.Println(string(json))

// 	session.Value = string(json)

// 	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).Update(st.sessionTableName).Set(session).ToSQL()

// 	if st.debug {
// 		log.Println(sqlStr)
// 	}

// 	_, errExec := st.db.Exec(sqlStr)

// 	if errExec != nil {
// 		return false, errExec
// 	}

// 	return true, nil
// }

// SQLCreateTable returns a SQL string for creating the cache table
func (st *Store) SQLCreateTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.sessionTableName + ` (
	  id varchar(40) NOT NULL PRIMARY KEY,
	  session_key varchar(40) NOT NULL,
	  session_value text,
	  expires_at datetime,
	  created_at datetime NOT NULL,
	  updated_at datetime NOT NULL,
	  deleted_at datetime
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.sessionTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "session_key" varchar(40) NOT NULL,
	  "session_value" text,
	  "expires_at" timestamptz(6),
	  "created_at" timestamptz(6) NOT NULL,
	  "updated_at" timestamptz(6) NOT NULL,
	  "deleted_at" timestamptz(6)
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.sessionTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "session_key" varchar(40) NOT NULL,
	  "session_value" text,
	  "expires_at" datetime,
	  "created_at" datetime NOT NULL,
	  "updated_at" datetime NOT NULL,
	  "deleted_at" datetime
	)
	`

	sql := "unsupported driver " + st.dbDriverName

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}
