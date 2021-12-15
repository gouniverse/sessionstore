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

func TestSessionDelete(t *testing.T) {
	db := InitDB("test_session_delete.db")

	store, _ := NewStore(WithDb(db), WithTableName("session"), WithAutoMigrate(true))

	sessionKey := "SESSION_KEY_DELETE"

	store.Set(sessionKey, "123456", 5)

	isDeleted, err := store.Delete(sessionKey)

	if err != nil {
		t.Fatalf("Session could not be deleted: " + err.Error())
	}

	if !isDeleted {
		t.Fatalf("Session remove key should return true on success")
	}

	if store.FindByKey(sessionKey) != nil {
		t.Fatalf("Session should no longer be present")
	}
}

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

func TestSetGetJSOM(t *testing.T) {
	db := InitDB("test_session_json.db")

	store, _ := NewStore(WithDb(db), WithTableName("session_json"), WithAutoMigrate(true))

	value := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	isSaved, err := store.SetJSON("mykey", value, 5)

	if err != nil {
		t.Fatalf("Set JSON failed: " + err.Error())
	}

	if !isSaved {
		t.Fatalf("Set JSON failed: " + err.Error())
	}

	result, err := store.GetJSON("mykey", "{}")

	if err != nil {
		t.Fatalf("Get JSON failed: " + err.Error())
	}

	var res = map[string]string{}
	for k, v := range result.(map[string]interface{}) {
		res[k] = v.(string)
	}

	if res["key1"] != value["key1"] {
		t.Fatalf("Key1 not correct: " + res["key1"])
	}

	if res["key2"] != value["key2"] {
		t.Fatalf("Key2 not correct: " + res["key2"])
	}

	if res["key3"] != value["key1"] {
		t.Fatalf("Key3 not correct: " + res["key3"])
	}
}
