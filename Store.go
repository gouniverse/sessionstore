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
	debugEnabled       bool
}

// NewStoreOptions define the options for creating a new session store
type NewStoreOptions struct {
	SessionTableName   string
	DB                 *sql.DB
	DbDriverName       string
	TimeoutSeconds     int64
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new session store
func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
		sessionTableName:   opts.SessionTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.sessionTableName == "" {
		return nil, errors.New("session store: sessionTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("session store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = store.DriverName(store.db)
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
	st.debugEnabled = debug
}

// ExpireSessionGoroutine - soft deletes expired cache
func (st *Store) ExpireSessionGoroutine() error {
	i := 0
	for {
		i++
		if st.debugEnabled {
			log.Println("Cleaning expired sessions...")
		}
		sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.sessionTableName).Where(goqu.C("expires_at").Lt(time.Now())).Delete().ToSQL()

		if st.debugEnabled {
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

	if st.debugEnabled {
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
func (st *Store) FindByKey(sessionKey string) (*Session, error) {
	// key exists, expires is < now, deleted null
	sqlStr, _, sqlErr := goqu.Dialect(st.dbDriverName).
		From(st.sessionTableName).
		Where(goqu.C("session_key").Eq(sessionKey), goqu.C("expires_at").Gt(time.Now()), goqu.C("deleted_at").IsNull()).
		Select("*").
		ToSQL()

	if sqlErr != nil {
		return nil, sqlErr
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	var session Session
	err := sqlscan.Get(context.Background(), st.db, &session, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// Looks like this is now outdated for sqlscan
			return nil, nil
		}
		if sqlscan.NotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &session, nil
}

// Gets the session value as a string
func (st *Store) Get(sessionKey string, valueDefault string) (string, error) {
	session, errFindByKey := st.FindByKey(sessionKey)

	if errFindByKey != nil {
		return "", errFindByKey
	}

	if session != nil {
		return session.Value, nil
	}

	return valueDefault, nil
}

// GetJSON attempts to parse the value as JSON, use with SetJSON
func (st *Store) GetJSON(key string, valueDefault interface{}) (interface{}, error) {
	session, errFindByKey := st.FindByKey(key)

	if errFindByKey != nil {
		var empty interface{}
		return empty, errFindByKey
	}

	if session != nil {
		jsonValue := session.Value
		var val interface{}
		jsonError := json.Unmarshal([]byte(jsonValue), &val)
		if jsonError != nil {
			return valueDefault, jsonError
		}

		return val, nil
	}

	return valueDefault, nil
}

// Has finds if a session by key exists
func (st *Store) Has(sessionKey string) (bool, error) {
	// key exists, expires is < now, deleted null
	sqlStr, _, sqlErr := goqu.Dialect(st.dbDriverName).
		From(st.sessionTableName).
		Where(goqu.C("session_key").Eq(sessionKey), goqu.C("expires_at").Gt(time.Now()), goqu.C("deleted_at").IsNull()).
		Select(goqu.COUNT("*")).As("count").
		ToSQL()

	if sqlErr != nil {
		return false, sqlErr
	}

	if st.debugEnabled {
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

		if st.debugEnabled {
			log.Println("SessionStore. Error: ", err)
		}
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// Set sets a key in store
func (st *Store) Set(sessionKey string, value string, seconds int64) error {
	session, errFindByKey := st.FindByKey(sessionKey)

	if errFindByKey != nil {
		return errFindByKey
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(seconds))

	var sqlStr string
	var sqlErr error
	if session == nil {
		var newSession = Session{
			ID:        uid.MicroUid(),
			Key:       sessionKey,
			Value:     value,
			ExpiresAt: &expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		sqlStr, _, sqlErr = goqu.Dialect(st.dbDriverName).
			Insert(st.sessionTableName).
			Rows(newSession).
			ToSQL()
	} else {
		fields := map[string]interface{}{}
		fields["session_value"] = value
		fields["expires_at"] = &expiresAt
		fields["updated_at"] = time.Now()
		sqlStr, _, sqlErr = goqu.Dialect(st.dbDriverName).
			Update(st.sessionTableName).
			Where(goqu.C("session_key").Eq(sessionKey), goqu.C("expires_at").Gt(time.Now()), goqu.C("deleted_at").IsNull()).
			Set(fields).
			ToSQL()
	}

	if sqlErr != nil {
		return sqlErr
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}

// SetJSON convenience method which saves the supplied value as JSON, use GetJSON to extract
func (st *Store) SetJSON(key string, value interface{}, seconds int64) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(key, string(jsonValue), seconds)
}
