package sessionstore

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/emirpasic/gods/maps/hashmap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Store defines a session store
type Store struct {
	sessionTableName string
	db               *gorm.DB
	timeoutSeconds   int
	automigrateEnabled bool
}

// StoreOption options for the session store
type StoreOption func(*Store)

// WithAutoMigrate sets the table name for the cache store
func WithAutoMigrate(automigrateEnabled bool) StoreOption {
	return func(s *Store) {
		s.automigrateEnabled = automigrateEnabled
	}
}

// WithDriverAndDNS sets the driver and the DNS for the database for the cache store
func WithDriverAndDNS(driverName string, dsn string) StoreOption {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return func(s *Store) {
		s.db = db
	}
}

// WithGormDb sets the GORM database for the session store
func WithGormDb(db *gorm.DB) StoreOption {
	return func(s *Store) {
		s.db = db
	}
}

// WithTableName sets the table name for the session store
func WithTableName(sessionTableName string) StoreOption {
	return func(s *Store) {
		s.sessionTableName = sessionTableName
	}
}

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) *Store {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.sessionTableName == "" {
		log.Panic("Session store: sessionTableName is required")
	}

	store.timeoutSeconds = 2 * 60 * 60 // 2 hours

	if store.automigrateEnabled == true {
		store.AutoMigrate()
	}

	return store
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() {
	st.db.Table(st.logTableName).AutoMigrate(&Session{})
}

// ExpireSessionGoroutine - soft deletes expired cache
func (st *Store) ExpireSessionGoroutine() {
	i := 0
	for {
		i++
		fmt.Println("Cleaning expired cache...")
		st.db.Table(st.sessionTableName).Where("`expires_at` < ?", time.Now()).Delete(Session{})
		time.Sleep(60 * time.Second) // Every minute
	}
}

// // SessionDelete removes all keys from the sessiom
// func SessionDelete(sessionKey string) bool {
// 	session := SessionFindByToken(sessionKey)

// 	if session == nil {
// 		return true
// 	}

// 	GetDb().Delete(&session)

// 	return true
// }

// FindByToken finds a session by key
func (st *Store) FindByToken(key string) *Session {
	session := &Session{}
	result := st.db.Table(st.sessionTableName).Where("`session_key` = ?", key).First(&session)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return session
}

// Get a key from session
func (st *Store) Get(key string, valueDefault string) string {
	cache := st.FindByToken(key)

	if cache != nil {
		return cache.Value
	}

	return valueDefault
}

// Start starts a session with a specified key
func (st *Store) Start(key string) bool {
	session := st.FindByToken(key)
	expiresAt := time.Now().Add(time.Duration(st.timeoutSeconds) * time.Second)

	if session != nil {
		return true
	}

	var newSession = Session{Key: key, Value: "{}", ExpiresAt: &expiresAt}

	dbResult := st.db.Table(st.sessionTableName).Create(&newSession)

	if dbResult.Error != nil {
		return false
	}

	return true
}

// Set sets a key in store
func (st *Store) Set(key string, value string, seconds int64) bool {
	session := st.FindByToken(key)
	expiresAt := time.Now().Add(time.Duration(st.timeoutSeconds) * time.Duration(seconds))

	if session != nil {
		session.Value = value
		session.ExpiresAt = &expiresAt
		//dbResult := GetDb().Table(User).Where("`key` = ?", key).Update(&cache)
		dbResult := st.db.Table(st.sessionTableName).Save(&session)
		if dbResult != nil {
			return false
		}
		return true
	}

	var newSessiom = Session{Key: key, Value: value, ExpiresAt: &expiresAt}

	dbResult := st.db.Table(st.sessionTableName).Create(&newSessiom)

	if dbResult.Error != nil {
		return false
	}

	return true
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
func (st *Store) Empty(sessionKey string) bool {
	session := st.FindByToken(sessionKey)

	kv := hashmap.New()

	if session == nil {
		return true
	}

	json, err := kv.ToJSON()

	if err != nil {
		return false
	}

	session.Value = string(json)

	st.db.Table(st.sessionTableName).Save(&session)

	return true
}

// SetKey gets a key from sessiom
func (st *Store) SetKey(sessionKey string, key string, value string) bool {
	session := st.FindByToken(sessionKey)

	kv := hashmap.New()

	if session == nil {
		isOk := st.Set(sessionKey, "{}", 2000)
		if isOk == false {
			return false
		}
		session = st.FindByToken(sessionKey)
	}

	log.Println(value)

	err := kv.FromJSON([]byte(session.Value))
	if err != nil {
		return false
	}

	kv.Put(key, value)
	json, err := kv.ToJSON()

	if err != nil {
		return false
	}
	log.Println(string(json))

	session.Value = string(json)

	st.db.Table(st.sessionTableName).Save(&session)

	return true
}

// RemoveKey removes a key from sessiom
func (st *Store) RemoveKey(sessionKey string, key string) bool {
	session := st.FindByToken(sessionKey)

	kv := hashmap.New()

	if session == nil {
		return true
	}

	err := kv.FromJSON([]byte(session.Value))
	if err != nil {
		return false
	}

	kv.Remove(key)

	json, err := kv.ToJSON()

	if err != nil {
		return false
	}

	log.Println(string(json))

	session.Value = string(json)

	st.db.Table(st.sessionTableName).Save(&session)

	return true
}
