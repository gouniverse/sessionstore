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

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session_create",
		AutomigrateEnabled: true,
		//DebugEnabled:       true,
	})

	if err != nil {
		t.Fatalf("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatalf("Store could not be created")
	}

	errSet := store.Set("post", "1234567890", 5)

	if errSet != nil {
		t.Fatalf("Cache could not be created: " + err.Error())
	}
}

func TestStoreAutomigrate(t *testing.T) {
	db := InitDB("test_session_automigrate.db")

	store, _ := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session_automigrate",
		AutomigrateEnabled: true,
	})

	err := store.AutoMigrate()

	if err != nil {
		t.Fatalf("Automigrate failed: " + err.Error())
	}

	errSet := store.Set("post", "1234567890", 5)

	if errSet != nil {
		t.Fatalf("Session could not be created: " + err.Error())
	}
}

func TestSessionDelete(t *testing.T) {
	db := InitDB("test_session_delete.db")

	store, _ := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session_delete",
		AutomigrateEnabled: true,
	})

	sessionKey := "SESSION_KEY_DELETE"

	store.Set(sessionKey, "123456", 5)

	isDeleted, err := store.Delete(sessionKey)

	if err != nil {
		t.Fatalf("Session could not be deleted: " + err.Error())
	}

	if !isDeleted {
		t.Fatalf("Session remove key should return true on success")
	}

	session, errFind := store.FindByKey(sessionKey)

	if errFind != nil {
		t.Fatal("Session find error", errFind.Error())
	}

	if session != nil {
		t.Fatal("Session should no longer be present")
	}
}

func TestStoreEnableDebug(t *testing.T) {
	db := InitDB("test_session_debug.db")

	store, _ := NewStore(NewStoreOptions{
		DB:               db,
		SessionTableName: "session_debug",
	})
	store.EnableDebug(true)

	err := store.AutoMigrate()

	if err != nil {
		t.Fatalf("Automigrate failed: " + err.Error())
	}
}

func TestSetKey(t *testing.T) {
	db := InitDB("test_session_set_key.db")

	store, _ := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session_set_key",
		AutomigrateEnabled: true,
	})

	err := store.Set("hello", "world", 1)

	if err != nil {
		t.Fatal("Setting key failed:", err.Error())
	}

	value, errGet := store.Get("hello", "")

	if errGet != nil {
		t.Fatal("Getting key failed: ", errGet.Error())
	}

	if value != "world" {
		t.Fatal("Incorrect value:", err.Error())
	}
}

func TestUpdateKey(t *testing.T) {
	db := InitDB("test_session_update_key.db")

	store, _ := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session_update_key",
		AutomigrateEnabled: true,
	})

	err := store.Set("hello", "world", 1)

	if err != nil {
		t.Fatalf("Setting key failed: " + err.Error())
	}

	session1, errFind := store.FindByKey("hello")

	if errFind != nil {
		t.Fatal("Session find error", errFind.Error())
	}

	if session1 == nil {
		t.Fatal("Session 1 not found")
	}

	time.Sleep(2 * time.Second)

	err2 := store.Set("hello", "world", 1)

	if err2 != nil {
		t.Fatalf("Update setting key failed: " + err2.Error())
	}

	session2, errFind := store.FindByKey("hello")

	if errFind != nil {
		t.Fatal("Session find error", errFind.Error())
	}

	if session2 == nil {
		t.Fatalf("Session 2 not found")
	}

	if session2.Value != "world" {
		t.Fatal("Value not correct:", session2.Value)
	}

	if session2.Key != "hello" {
		t.Fatal("Key not correct:", session2.Key)
	}

	if session2.UpdatedAt == session1.CreatedAt {
		t.Fatal("Updated at should be different from created at date:", session2.UpdatedAt.Format(time.UnixDate))
	}

	if session2.UpdatedAt.Sub(session1.CreatedAt).Seconds() < 1 {
		t.Fatal("Updated at should more than 1 second after created at date:", session2.UpdatedAt.Format(time.UnixDate), " - ", session1.CreatedAt.Format(time.UnixDate))
	}
}

func TestSetGetJSOM(t *testing.T) {
	db := InitDB("test_session_json.db")

	store, _ := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session_json",
		AutomigrateEnabled: true,
	})

	value := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	err := store.SetJSON("mykey", value, 5)

	if err != nil {
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

	if res["key3"] != value["key3"] {
		t.Fatalf("Key3 not correct: " + res["key3"])
	}
}

func TestHas(t *testing.T) {
	db := InitDB("test_session_has.db")

	store, _ := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session_has",
		AutomigrateEnabled: true,
	})

	hasNo, err := store.Has("mykey")

	if err != nil {
		t.Fatalf("Has no failed: " + err.Error())
	}

	if hasNo {
		t.Fatalf("Has no failed: " + err.Error())
	}

	errSet := store.Set("mykey", "test", 5)

	if errSet != nil {
		t.Fatalf("Set failed: " + err.Error())
	}

	has, err := store.Has("mykey")

	if err != nil {
		t.Fatalf("Has failed: " + err.Error())
	}

	if !has {
		t.Fatalf("Has failed: " + err.Error())
	}
}
