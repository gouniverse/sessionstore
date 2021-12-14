package sessionstore

import (
	"database/sql"
	"os"
	"testing"
	"time"

	// "time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("test_session_store_create.db")

	store, err := NewStore(WithDb(db), WithTableName("session_create"), WithAutoMigrate(true))

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	isOk, err := store.Set("post", "1234567890", 5)

	if err != nil {
		t.Fatalf("Cache could not be created: " + err.Error())
	}

	if isOk == false {
		t.Fatalf("Cache could not be created")
	}
}

func TestStoreAutomigrate(t *testing.T) {
	db := InitDB("test_session_automigrate.db")

	store, _ := NewStore(WithDb(db), WithTableName("session_automigrate"))

	err := store.AutoMigrate()

	if err != nil {
		t.Fatalf("Automigrate failed: " + err.Error())
	}

	isOk, err := store.Set("post", "1234567890", 5)

	if err != nil {
		t.Fatalf("Session could not be created: " + err.Error())
	}

	if isOk == false {
		t.Fatalf("Session could not be created")
	}
}

// func TestStoreSessionKeyRemove(t *testing.T) {
// 	db := InitDB("test_session_delete.db")

// 	store, _ := NewStore(WithDb(db), WithTableName("session"), WithAutoMigrate(true))

// 	isRemoved, err := store.RemoveKey("post", "key1")

// 	if err != nil {
// 		t.Fatalf("Session could not be removed: " + err.Error())
// 	}

// 	if !isRemoved {
// 		t.Fatalf("Session remove key should return true on success")
// 	}

// 	if store.FindBySessionKey("post") != nil {
// 		t.Fatalf("Cache should no longer be present")
// 	}
// }

func TestStoreEnableDebug(t *testing.T) {
	db := InitDB("test_session_debug.db")

	store, _ := NewStore(WithDb(db), WithTableName("session_debug"))
	store.EnableDebug(true)

	err := store.AutoMigrate()

	if err != nil {
		t.Fatalf("Automigrate failed: " + err.Error())
	}
}

func TestSetKey(t *testing.T) {
	db := InitDB("test_session_set_key.db")

	store, _ := NewStore(WithDb(db), WithTableName("session_key"), WithAutoMigrate(true))

	ok, err := store.Set("hello", "world", 1)

	if err != nil {
		t.Fatalf("Setting key failed: " + err.Error())
	}

	if ok != true {
		t.Fatalf("Response not true: " + err.Error())
	}

	value := store.Get("hello", "")

	if value != "world" {
		t.Fatalf("Incorrect value: " + err.Error())
	}
}

func TestUpdateKey(t *testing.T) {
	db := InitDB("test_session_update_key.db")

	store, _ := NewStore(WithDb(db), WithTableName("session_update"), WithAutoMigrate(true))

	ok, err := store.Set("hello", "world", 1)

	if err != nil {
		t.Fatalf("Setting key failed: " + err.Error())
	}

	if ok != true {
		t.Fatalf("Response not true: " + err.Error())
	}

	session1 := store.FindByKey("hello")

	time.Sleep(2 * time.Second)

	ok2, err2 := store.Set("hello", "world", 1)

	if err2 != nil {
		t.Fatalf("Update setting key failed: " + err2.Error())
	}

	if ok2 != true {
		t.Fatalf("Update response not true: " + err.Error())
	}

	session2 := store.FindByKey("hello")

	if session2 == nil {
		t.Fatalf("Cache not found: " + err.Error())
	}

	if session2.Value != "world" {
		t.Fatalf("Value not correct: " + session2.Value)
	}

	if session2.Key != "hello" {
		t.Fatalf("Key not correct: " + session2.Key)
	}

	if session2.UpdatedAt == session1.CreatedAt {
		t.Fatalf("Updated at should be different from created at date: " + session2.UpdatedAt.Format(time.UnixDate))
	}

	if session2.UpdatedAt.Sub(session1.CreatedAt).Seconds() < 1 {
		t.Fatalf("Updated at should more than 1 second after created at date: " + session2.UpdatedAt.Format(time.UnixDate) + " - " + session1.CreatedAt.Format(time.UnixDate))
	}
}
