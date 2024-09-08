package sessionstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"     // importing mysql dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"  // importing postgres dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"   // importing sqlite3 dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver" // importing sqlserver dialect
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

var _ StoreInterface = (*Store)(nil) // verify it extends the task interface

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
		store.dbDriverName = sb.DatabaseDriverName(store.db)
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
		sqlStr, sqlParams, err := goqu.Dialect(st.dbDriverName).
			From(st.sessionTableName).
			Where(goqu.C(COLUMN_EXPIRES_AT).Lt(time.Now())).
			Delete().
			Prepared(true).
			ToSQL()

		if err != nil {
			return err
		}

		if st.debugEnabled {
			log.Println(sqlStr)
		}

		_, err = st.db.Exec(sqlStr, sqlParams...)

		if err != nil {
			if err == sql.ErrNoRows {
				// Looks like this is now outdated for sqlscan
				return nil
			}

			if sqlscan.NotFound(err) {
				return nil
			}

			log.Println("Session Store. ExpireSessionGoroutine. Error: ", err)
			return nil
		}

		time.Sleep(60 * time.Second) // Every minute
	}
}

func (st *Store) Extend(sessionKey string, seconds int64, options SessionOptions) error {
	session, errFindByKey := st.FindByKey(sessionKey, options)

	if errFindByKey != nil {
		return errFindByKey
	}

	if session == nil {
		return errors.New("session not found")
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(seconds))

	session.ExpiresAt = &expiresAt

	err := st.sessionUpdate(*session, options)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a session
func (st *Store) Delete(sessionKey string, options SessionOptions) (bool, error) {
	wheres := []goqu.Expression{
		goqu.C(COLUMN_SESSION_KEY).Eq(sessionKey),
		goqu.C(COLUMN_USER_AGENT).Eq(options.UserAgent),
		goqu.C(COLUMN_IP_ADDRESS).Eq(options.IPAddress),
	}

	// Only add the condition, if specifically requested
	if len(options.UserID) > 0 {
		wheres = append(wheres, goqu.C(COLUMN_USER_ID).Eq(options.UserID))
	}

	sqlStr, sqlParams, err := goqu.Dialect(st.dbDriverName).
		From(st.sessionTableName).
		Where(wheres...).
		Delete().
		Prepared(true).
		ToSQL()

	if err != nil {
		return false, err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err = st.db.Exec(sqlStr, sqlParams...)

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
func (st *Store) FindByKey(sessionKey string, options SessionOptions) (*Session, error) {
	wheres := []goqu.Expression{
		goqu.C(COLUMN_SESSION_KEY).Eq(sessionKey),
		goqu.C(COLUMN_EXPIRES_AT).Gt(time.Now()),
		goqu.C(COLUMN_DELETED_AT).Eq(time.Time{}),
		goqu.C(COLUMN_USER_AGENT).Eq(options.UserAgent),
		goqu.C(COLUMN_IP_ADDRESS).Eq(options.IPAddress),
	}

	// Only add the condition, if specifically requested
	if len(options.UserID) > 0 {
		wheres = append(wheres, goqu.C(COLUMN_USER_ID).Eq(options.UserID))
	}

	// key exists, expires is < now, deleted null
	sqlStr, _, sqlErr := goqu.Dialect(st.dbDriverName).
		From(st.sessionTableName).
		Where(wheres...).
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
func (st *Store) Get(sessionKey string, valueDefault string, options SessionOptions) (string, error) {
	session, errFindByKey := st.FindByKey(sessionKey, options)

	if errFindByKey != nil {
		return "", errFindByKey
	}

	if session != nil {
		return session.Value, nil
	}

	return valueDefault, nil
}

// GetAny attempts to parse the value as interface, use with SetAny
func (st *Store) GetAny(key string, valueDefault interface{}, options SessionOptions) (interface{}, error) {
	session, errFindByKey := st.FindByKey(key, options)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
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

// GetMap attempts to parse the value as map[string]any, use with SetMap
func (st *Store) GetMap(key string, valueDefault map[string]any, options SessionOptions) (map[string]any, error) {
	session, errFindByKey := st.FindByKey(key, options)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
	}

	if session != nil {
		jsonValue := session.Value
		var val map[string]any
		jsonError := json.Unmarshal([]byte(jsonValue), &val)
		if jsonError != nil {
			return valueDefault, jsonError
		}

		return val, nil
	}

	return valueDefault, nil
}

// Has finds if a session by key exists
func (st *Store) Has(sessionKey string, options SessionOptions) (bool, error) {
	wheres := []goqu.Expression{
		goqu.C(COLUMN_SESSION_KEY).Eq(sessionKey),
		goqu.C(COLUMN_EXPIRES_AT).Gt(time.Now()),
		goqu.C(COLUMN_DELETED_AT).Eq(time.Time{}),
		goqu.C(COLUMN_USER_AGENT).Eq(options.UserAgent),
		goqu.C(COLUMN_IP_ADDRESS).Eq(options.IPAddress),
	}

	// Only add the condition, if specifically requested
	if len(options.UserID) > 0 {
		wheres = append(wheres, goqu.C(COLUMN_USER_ID).Eq(options.UserID))
	}

	// key exists, expires is < now, deleted null
	sqlStr, _, sqlErr := goqu.Dialect(st.dbDriverName).
		From(st.sessionTableName).
		Where(wheres...).
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

func (st *Store) MergeMap(key string, mergeMap map[string]any, seconds int64, options SessionOptions) error {
	currentMap, err := st.GetMap(key, nil, options)

	if err != nil {
		return err
	}

	if currentMap == nil {
		return errors.New("sessionstore. nil found")
	}

	for mapKey, mapValue := range mergeMap {
		currentMap[mapKey] = mapValue
	}

	return st.SetMap(key, currentMap, seconds, options)
}

func (st *Store) sessionCreate(session Session) error {
	sqlStr, _, sqlErr := goqu.Dialect(st.dbDriverName).
		Insert(st.sessionTableName).
		Rows(session).
		ToSQL()

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

func (st *Store) sessionUpdate(session Session, options SessionOptions) error {
	fields := map[string]interface{}{}
	fields[COLUMN_SESSION_VALUE] = session.Value
	fields[COLUMN_EXPIRES_AT] = session.ExpiresAt
	fields[COLUMN_UPDATED_AT] = time.Now()

	wheres := []goqu.Expression{
		goqu.C(COLUMN_SESSION_KEY).Eq(session.Key),
		goqu.C(COLUMN_EXPIRES_AT).Gt(time.Now()),
		goqu.C(COLUMN_DELETED_AT).Eq(time.Time{}),
		goqu.C(COLUMN_USER_AGENT).Eq(options.UserAgent),
		goqu.C(COLUMN_IP_ADDRESS).Eq(options.IPAddress),
	}

	// Only add the condition, if specifically requested
	if len(options.UserID) > 0 {
		wheres = append(wheres, goqu.C(COLUMN_USER_ID).Eq(options.UserID))
	}

	sqlStr, _, sqlErr := goqu.Dialect(st.dbDriverName).
		Update(st.sessionTableName).
		Where(wheres...).
		Set(fields).
		ToSQL()

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

// Set sets a key in store
func (st *Store) Set(sessionKey string, value string, seconds int64, options SessionOptions) error {
	session, errFindByKey := st.FindByKey(sessionKey, options)

	if errFindByKey != nil {
		return errFindByKey
	}

	expiresAt := time.Now().Add(time.Second * time.Duration(seconds))

	if session == nil {
		var newSession = Session{
			ID:        uid.MicroUid(),
			Key:       sessionKey,
			Value:     value,
			UserID:    options.UserID,
			UserAgent: options.UserAgent,
			IPAddress: options.IPAddress,
			ExpiresAt: &expiresAt,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: &time.Time{},
		}

		return st.sessionCreate(newSession)
	} else {
		session.Value = value
		session.ExpiresAt = &expiresAt
		session.UpdatedAt = time.Now()

		return st.sessionUpdate(*session, options)
	}
}

// SetAny convenience method which saves the supplied interface value, use GetAny to extract
// Internally it serializes the data to JSON
func (st *Store) SetAny(key string, value interface{}, seconds int64, options SessionOptions) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(key, string(jsonValue), seconds, options)
}

// SetMap convenience method which saves the supplied map, use GetMap to extract
func (st *Store) SetMap(key string, value map[string]any, seconds int64, options SessionOptions) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(key, string(jsonValue), seconds, options)
}
