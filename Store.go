package sessionstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"     // importing mysql dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"  // importing postgres dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"   // importing sqlite3 dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver" // importing sqlserver dialect
	"github.com/dromara/carbon/v2"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// == INTERFACE ===============================================================

var _ StoreInterface = (*store)(nil) // verify it extends the task interface

// == TYPE ====================================================================

// Store defines a session store
type store struct {
	sessionTableName   string
	db                 *sql.DB
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool
	sqlLogger          *slog.Logger
}

// PUBLIC METHODS ============================================================

// AutoMigrate auto migrate
func (store *store) AutoMigrate() error {
	sqlStr := store.SQLCreateTable()

	if sqlStr == "" {
		return errors.New("session store: table create sql is empty")
	}

	if store.db == nil {
		return errors.New("session store: database is nil")
	}

	_, err := store.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// ExpireSessionGoroutine - soft deletes expired sessions
func (st *store) ExpireSessionGoroutine() error {
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

func (store *store) Extend(sessionKey string, seconds int64, options SessionOptions) error {
	session, errFindByKey := store.FindByKey(sessionKey, options)

	if errFindByKey != nil {
		return errFindByKey
	}

	if session == nil {
		return errors.New("session not found")
	}

	expiresAt := carbon.Now(carbon.UTC).AddSeconds(cast.ToInt(seconds)).ToDateTimeString(carbon.UTC)

	session.SetExpiresAt(expiresAt)

	err := store.SessionUpdate(session)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a session
func (st *store) Delete(sessionKey string, options SessionOptions) error {
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

		return err
	}

	return nil
}

// FindByKey finds a session by key
func (store *store) FindByKey(sessionKey string, options SessionOptions) (SessionInterface, error) {
	if sessionKey == "" {
		return nil, errors.New("session store > find by key: session key is required")
	}

	query := SessionQuery().
		SetKey(sessionKey).
		SetExpiresAtGte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUserAgent(options.UserAgent).
		SetUserIpAddress(options.IPAddress).
		SetLimit(1)

	// Only add the UserID, if specifically requested
	if len(options.UserID) > 0 {
		query.SetUserID(options.UserID)
	}

	list, err := store.SessionList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// Gets the session value as a string
func (st *store) Get(sessionKey string, valueDefault string, options SessionOptions) (string, error) {
	session, errFindByKey := st.FindByKey(sessionKey, options)

	if errFindByKey != nil {
		return "", errFindByKey
	}

	if session != nil {
		return session.GetValue(), nil
	}

	return valueDefault, nil
}

// GetAny attempts to parse the value as interface, use with SetAny
func (st *store) GetAny(key string, valueDefault interface{}, options SessionOptions) (interface{}, error) {
	session, errFindByKey := st.FindByKey(key, options)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
	}

	if session != nil {
		jsonValue := session.GetValue()
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
func (st *store) GetMap(key string, valueDefault map[string]any, options SessionOptions) (map[string]any, error) {
	session, errFindByKey := st.FindByKey(key, options)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
	}

	if session != nil {
		jsonValue := session.GetValue()
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
func (store *store) Has(sessionKey string, options SessionOptions) (bool, error) {
	if sessionKey == "" {
		return false, errors.New("session store > find by key: session key is required")
	}

	query := SessionQuery().
		SetKey(sessionKey).
		SetExpiresAtGte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUserAgent(options.UserAgent).
		SetUserIpAddress(options.IPAddress).
		SetLimit(1)

	// Only add the UserID, if specifically requested
	if len(options.UserID) > 0 {
		query.SetUserID(options.UserID)
	}

	count, err := store.SessionCount(query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (st *store) MergeMap(key string, mergeMap map[string]any, seconds int64, options SessionOptions) error {
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

func (store *store) SessionCount(options SessionQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.sessionSelectQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	mapped, err := db.SelectToMapString(sqlStr, params...)
	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (st *store) SessionCreate(session SessionInterface) error {
	if session == nil {
		return errors.New("sessionstore > session create. session cannot be nil")
	}

	if session.GetKey() == "" {
		return errors.New("sessionstore > session create. key cannot be empty")
	}

	if session.GetExpiresAt() == "" {
		return errors.New("sessionstore > session create. expires at cannot be empty")
	}

	if session.GetCreatedAt() == "" {
		session.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}

	if session.GetUpdatedAt() == "" {
		session.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}

	if session.GetSoftDeletedAt() == "" {
		session.SetSoftDeletedAt(sb.MAX_DATETIME)
	}

	data := session.Data()

	sqlStr, sqlParams, sqlErr := goqu.Dialect(st.dbDriverName).
		Insert(st.sessionTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if sqlErr != nil {
		return sqlErr
	}

	st.logSql("create", sqlStr, sqlParams)

	_, err := st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	session.MarkAsNotDirty()

	return nil
}

// SessionDelete deletes a session
func (store *store) SessionDelete(session SessionInterface) error {
	if session == nil {
		return errors.New("session is nil")
	}

	return store.SessionDeleteByID(session.GetID())
}

// SessionDeleteByID deletes a session by id
func (store *store) SessionDeleteByID(id string) error {
	if id == "" {
		return errors.New("session id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.sessionTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

// SessionDeleteByID deletes a session by id
func (store *store) SessionDeleteByKey(sessionKey string) error {
	if sessionKey == "" {
		return errors.New("session id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.sessionTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_SESSION_KEY).Eq(sessionKey)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

// SessionFindByID finds a session by id
func (store *store) SessionFindByID(sessionID string) (SessionInterface, error) {
	if sessionID == "" {
		return nil, errors.New("session store > find by id: session id is required")
	}

	query := SessionQuery().
		SetID(sessionID).
		SetExpiresAtGte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetLimit(1)

	list, err := store.SessionList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// SessionFindByKey finds a session by key
func (store *store) SessionFindByKey(sessionKey string) (SessionInterface, error) {
	if sessionKey == "" {
		return nil, errors.New("session store > find by key: session key is required")
	}

	query := SessionQuery().
		SetKey(sessionKey).
		SetExpiresAtGte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetLimit(1)

	list, err := store.SessionList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) SessionList(query SessionQueryInterface) ([]SessionInterface, error) {
	if query == nil {
		return []SessionInterface{}, errors.New("at session list > session query is nil")
	}

	q, columns, err := store.sessionSelectQuery(query)

	if err != nil {
		return []SessionInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []SessionInterface{}, nil
	}

	store.logSql("list", sqlStr, sqlParams...)

	if store.db == nil {
		return []SessionInterface{}, errors.New("userstore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []SessionInterface{}, errors.New("userstore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr, sqlParams...)

	if err != nil {
		return []SessionInterface{}, err
	}

	list := []SessionInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewSessionFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) SessionSoftDelete(session SessionInterface) error {
	if session == nil {
		return errors.New("session is nil")
	}

	session.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SessionUpdate(session)
}

func (store *store) SessionSoftDeleteByID(id string) error {
	session, err := store.SessionFindByID(id)

	if err != nil {
		return err
	}

	return store.SessionSoftDelete(session)
}

func (store *store) SessionUpdate(session SessionInterface) error {
	if session == nil {
		return errors.New("sessionstore > session update. session cannot be nil")
	}

	if store.db == nil {
		return errors.New("sessionstore > session update. db cannot be nil")
	}

	session.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	dataChanged := session.DataChanged()

	if len(dataChanged) == 0 {
		return nil
	}

	delete(dataChanged, COLUMN_ID) // ID cannot be updated

	// fields := map[string]interface{}{}
	// fields[COLUMN_SESSION_VALUE] = session.GetValue()
	// fields[COLUMN_EXPIRES_AT] = session.GetExpiresAt()
	// fields[COLUMN_UPDATED_AT] = carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)

	// wheres := []goqu.Expression{
	// 	goqu.C(COLUMN_SESSION_KEY).Eq(session.GetKey()),
	// 	goqu.C(COLUMN_EXPIRES_AT).Gte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)),
	// 	goqu.C(COLUMN_SOFT_DELETED_AT).Gte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)),
	// 	goqu.C(COLUMN_USER_AGENT).Eq(options.UserAgent),
	// 	goqu.C(COLUMN_IP_ADDRESS).Eq(options.IPAddress),
	// }

	// // Only add the condition, if specifically requested
	// if len(options.UserID) > 0 {
	// 	wheres = append(wheres, goqu.C(COLUMN_USER_ID).Eq(options.UserID))
	// }

	sqlStr, sqlParams, sqlErr := goqu.Dialect(store.dbDriverName).
		Update(store.sessionTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_SESSION_KEY).Eq(session.GetKey())).
		Where(goqu.C(COLUMN_ID).Eq(session.GetID())).
		Set(dataChanged).
		ToSQL()

	if sqlErr != nil {
		return sqlErr
	}

	store.logSql("update", sqlStr, sqlParams...)

	_, err := store.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	return nil
}

// Deprecated: Set sets a key in store
func (st *store) Set(sessionKey string, value string, seconds int64, options SessionOptions) error {
	session, errFindByKey := st.FindByKey(sessionKey, options)

	if errFindByKey != nil {
		return errFindByKey
	}

	expiresAt := carbon.Now(carbon.UTC).AddSeconds(cast.ToInt(seconds)).ToDateTimeString(carbon.UTC)

	if session == nil {
		newSession := NewSession().
			SetKey(sessionKey).
			SetValue(value).
			SetUserID(options.UserID).
			SetUserAgent(options.UserAgent).
			SetIPAddress(options.IPAddress).
			SetExpiresAt(expiresAt)

		return st.SessionCreate(newSession)
	} else {
		session.SetValue(value)
		session.SetExpiresAt(expiresAt)
		session.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

		return st.SessionUpdate(session)
	}
}

// Deprecated: SetAny convenience method which saves the supplied interface value, use GetAny to extract
// Internally it serializes the data to JSON
func (st *store) SetAny(key string, value interface{}, seconds int64, options SessionOptions) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(key, string(jsonValue), seconds, options)
}

// Deprecated: SetMap convenience method which saves the supplied map, use GetMap to extract
func (st *store) SetMap(key string, value map[string]any, seconds int64, options SessionOptions) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(key, string(jsonValue), seconds, options)
}

func (store *store) sessionSelectQuery(options SessionQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("session query: cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.sessionTableName)

	if options.HasCreatedAtGte() && options.HasCreatedAtLte() {
		q = q.Where(
			goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()),
			goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()),
		)
	} else if options.HasCreatedAtGte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Gte(options.CreatedAtGte()))
	} else if options.HasCreatedAtLte() {
		q = q.Where(goqu.C(COLUMN_CREATED_AT).Lte(options.CreatedAtLte()))
	}

	if options.HasExpiresAtGte() && options.HasExpiresAtLte() {
		q = q.Where(
			goqu.C(COLUMN_EXPIRES_AT).Gte(options.ExpiresAtGte()),
			goqu.C(COLUMN_EXPIRES_AT).Lte(options.ExpiresAtLte()),
		)
	} else if options.HasExpiresAtGte() {
		q = q.Where(goqu.C(COLUMN_EXPIRES_AT).Gte(options.ExpiresAtGte()))
	} else if options.HasExpiresAtLte() {
		q = q.Where(goqu.C(COLUMN_EXPIRES_AT).Lte(options.ExpiresAtLte()))
	}

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasKey() {
		q = q.Where(goqu.C(COLUMN_SESSION_KEY).Eq(options.Key()))
	}

	if options.HasUserAgent() {
		q = q.Where(goqu.C(COLUMN_USER_AGENT).Eq(options.UserAgent()))
	}

	if options.HasUserID() {
		q = q.Where(goqu.C(COLUMN_USER_ID).Eq(options.UserID()))
	}

	if options.HasUserIpAddress() {
		q = q.Where(goqu.C(COLUMN_IP_ADDRESS).Eq(options.UserIpAddress()))
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(uint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(uint(options.Offset()))
		}
	}

	sortOrder := sb.DESC
	if options.HasSortOrder() && options.SortOrder() != "" {
		sortOrder = options.SortOrder()
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted sessions requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}

func (store *store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}
