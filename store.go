package sessionstore

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Store defines a session store
type Store struct {
	sessionTableName string
	db               *gorm.DB
}

// StoreOption options for the session store
type StoreOption func(*Store)

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

	store.db.Table(store.sessionTableName).AutoMigrate(&Session{})

	return store
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
