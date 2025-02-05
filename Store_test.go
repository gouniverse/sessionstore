package sessionstore

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(filepath string) (*sql.DB, error) {
	if filepath != ":memory:" && utils.FileExists(filepath) {
		err := os.Remove(filepath) // remove database

		if err != nil {
			return nil, err
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initStore(filepath string) (StoreInterface, error) {
	db, err := initDB(filepath)

	if err != nil {
		return nil, err
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		SessionTableName:   "session",
		AutomigrateEnabled: true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	return store, nil
}

func TestStore_Create(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}
}

func TestStore_Automigrate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	err = store.AutoMigrate()

	if err != nil {
		t.Fatal("Automigrate failed: " + err.Error())
	}
}

// func TestSessionDelete(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	sessionKey := "SESSION_KEY_DELETE"

// 	store.Set(sessionKey, "123456", 5, SessionOptions{})

// 	err = store.Delete(sessionKey, SessionOptions{})

// 	if err != nil {
// 		t.Fatalf("Session could not be deleted: " + err.Error())
// 	}

// 	session, errFind := store.FindByKey(sessionKey, SessionOptions{})

// 	if errFind != nil {
// 		t.Fatal("Session find error", errFind.Error())
// 	}

// 	if session != nil {
// 		t.Fatal("Session should no longer be present")
// 	}
// }

func TestStore_EnableDebug(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	store.EnableDebug(true)

	err = store.AutoMigrate()

	if err != nil {
		t.Fatal("Automigrate failed: " + err.Error())
	}
}

// func TestSetKey(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	err = store.Set("hello", "world", 600, SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if err != nil {
// 		t.Fatal("Setting key failed:", err.Error())
// 	}

// 	value, errGet := store.Get("hello", "", SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if errGet != nil {
// 		t.Fatal("Getting key failed: ", errGet.Error())
// 	}

// 	if value != "world" {
// 		t.Fatal("Incorrect value:", value)
// 	}
// }

// func TestUpdateKey(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	err = store.Set("hello", "world", 1, SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if err != nil {
// 		t.Fatalf("Setting key failed: " + err.Error())
// 	}

// 	session1, errFind := store.FindByKey("hello", SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if errFind != nil {
// 		t.Fatal("Session find error", errFind.Error())
// 	}

// 	if session1 == nil {
// 		t.Fatal("Session 1 not found")
// 	}

// 	time.Sleep(2 * time.Second)

// 	err2 := store.Set("hello", "world", 1, SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if err2 != nil {
// 		t.Fatalf("Update setting key failed: " + err2.Error())
// 	}

// 	session2, errFind := store.FindByKey("hello", SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if errFind != nil {
// 		t.Fatal("Session find error", errFind.Error())
// 	}

// 	if session2 == nil {
// 		t.Fatalf("Session 2 not found")
// 	}

// 	if session2.GetValue() != "world" {
// 		t.Fatal("Value not correct:", session2.GetValue())
// 	}

// 	if session2.GetKey() != "hello" {
// 		t.Fatal("Key not correct:", session2.GetKey())
// 	}

// 	if session2.GetUpdatedAt() == session1.GetUpdatedAt() {
// 		t.Fatal("Updated at should be different from created at date:", session2.GetUpdatedAt())
// 	}

// 	if session2.GetUpdatedAtCarbon().DiffAbsInSeconds(session1.GetCreatedAtCarbon()) < 1 {
// 		t.Fatal("Updated at should more than 1 second after created at date:", session2.GetUpdatedAt(), " - ", session1.GetCreatedAt())
// 	}
// }

// func TestSetGetAny(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	value := map[string]string{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	err = store.SetAny("my.key", value, 5, SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if err != nil {
// 		t.Fatalf("Set JSON failed: " + err.Error())
// 	}

// 	result, err := store.GetAny("my.key", "{}", SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if err != nil {
// 		t.Fatalf("Get JSON failed: " + err.Error())
// 	}

// 	t.Log("Result: ", result)

// 	var res = map[string]string{}
// 	for k, v := range result.(map[string]interface{}) {
// 		res[k] = v.(string)
// 	}

// 	if res["key1"] != value["key1"] {
// 		t.Fatalf("Key1 not correct: " + res["key1"])
// 	}

// 	if res["key2"] != value["key2"] {
// 		t.Fatalf("Key2 not correct: " + res["key2"])
// 	}

// 	if res["key3"] != value["key3"] {
// 		t.Fatalf("Key3 not correct: " + res["key3"])
// 	}
// }

// func TestHas(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	hasNo, err := store.Has("mykey", SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if err != nil {
// 		t.Fatalf("Has no failed: " + err.Error())
// 	}

// 	if hasNo {
// 		t.Fatal("Has no failed: ", hasNo)
// 	}

// 	errSet := store.Set("mykey", "test", 5, SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if errSet != nil {
// 		t.Fatalf("Set failed: " + errSet.Error())
// 	}

// 	has, err := store.Has("mykey", SessionOptions{
// 		UserID:    "123456",
// 		UserAgent: "UserAgent",
// 		IPAddress: "127.0.0.1",
// 	})

// 	if err != nil {
// 		t.Fatalf("Has failed: " + err.Error())
// 	}

// 	if !has {
// 		t.Fatal("Has failed: ", has)
// 	}
// }

// func TestSetGetMap(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	value := map[string]any{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	err = store.SetMap("mykey", value, 5, SessionOptions{})

// 	if err != nil {
// 		t.Fatalf("Set Map failed: " + err.Error())
// 	}

// 	result, err := store.GetMap("mykey", nil, SessionOptions{})

// 	if err != nil {
// 		t.Fatalf("Get JSON failed: " + err.Error())
// 	}

// 	if result == nil {
// 		t.Fatalf("GetMap failed: nil returned")
// 	}

// 	if result["key1"].(string) != value["key1"] {
// 		t.Fatalf("Key1 not correct: " + result["key1"].(string))
// 	}

// 	if result["key2"] != value["key2"] {
// 		t.Fatalf("Key2 not correct: " + result["key2"].(string))
// 	}

// 	if result["key3"] != value["key3"] {
// 		t.Fatalf("Key3 not correct: " + result["key3"].(string))
// 	}
// }

// func TestMergeMap(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	value := map[string]any{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	err = store.SetMap("mykey", value, 600, SessionOptions{})

// 	if err != nil {
// 		t.Fatalf("Set Map failed: " + err.Error())
// 	}

// 	valueMerge := map[string]any{
// 		"key2": "value22",
// 		"key3": "value33",
// 	}

// 	err = store.MergeMap("mykey", valueMerge, 600, SessionOptions{})

// 	if err != nil {
// 		t.Fatalf("Merge Map failed: " + err.Error())
// 	}

// 	result, err := store.GetMap("mykey", nil, SessionOptions{})

// 	if err != nil {
// 		t.Fatalf("Get JSON failed: " + err.Error())
// 	}

// 	if result == nil {
// 		t.Fatalf("GetMap failed: nil returned")
// 	}

// 	if result["key1"].(string) != value["key1"] {
// 		t.Fatalf("Key1 not correct: " + result["key1"].(string))
// 	}

// 	if result["key2"].(string) != valueMerge["key2"] {
// 		t.Fatalf("Key2 not correct: " + result["key2"].(string))
// 	}

// 	if result["key3"].(string) != valueMerge["key3"] {
// 		t.Fatalf("Key3 not correct: " + result["key3"].(string))
// 	}
// }

// func TestExtend(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	err = store.Set("mykey", "test", 5, SessionOptions{})

// 	if err != nil {
// 		t.Fatal("Set failed: " + err.Error())
// 	}

// 	err = store.Extend("mykey", 100, SessionOptions{})

// 	if err != nil {
// 		t.Fatal("Extend failed: " + err.Error())
// 	}

// 	sessionExtended, err := store.FindByKey("mykey", SessionOptions{})

// 	if err != nil {
// 		t.Fatal("Extend failed: " + err.Error())
// 	}

// 	if sessionExtended == nil {
// 		t.Fatal("Extend failed. Session is NIL")
// 	}

// 	if sessionExtended.GetValue() != "test" {
// 		t.Fatal("Extend failed. Value is wrong", sessionExtended.GetValue())
// 	}

// 	diff := sessionExtended.GetExpiresAtCarbon().DiffAbsInSeconds(carbon.Now(carbon.UTC))

// 	if diff < 90 {
// 		t.Fatal("Extend failed. ExpiresAt must be more than 90 seconds", sessionExtended.GetExpiresAt(), diff)
// 	}

// 	if diff > 110 {
// 		t.Fatal("Extend failed. ExpiresAt must be less than 110 seconds", sessionExtended.GetExpiresAt(), diff)
// 	}

// }

func TestStore_SessionCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession()

	if session == nil {
		t.Fatal("unexpected nil session")
	}

	if session.GetID() == "" {
		t.Fatal("unexpected empty id:", session.GetID())
	}

	if len(session.GetID()) != 32 {
		t.Fatal("unexpected id length:", len(session.GetID()))
	}

	if session.GetKey() == "" {
		t.Fatal("unexpected empty key:", session.GetKey())
	}

	if len(session.GetKey()) != 100 {
		t.Fatal("unexpected key length:", len(session.GetKey()))
	}

	err = store.SessionCreate(session)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStore_SessionDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession()

	err = store.SessionCreate(session)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SessionDeleteByID(session.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sessionFindWithDeleted, err := store.SessionList(SessionQuery().
		SetID(session.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(sessionFindWithDeleted) != 0 {
		t.Fatal("Session MUST be deleted, but it is not")
	}
}

func TestStore_SessionDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession().
		SetValue("one two three four")

	if session == nil {
		t.Fatal("unexpected nil session")
	}

	if session.GetID() == "" {
		t.Fatal("unexpected empty id:", session.GetID())
	}

	err = store.SessionCreate(session)

	if err != nil {
		t.Error("unexpected error:", err)
	}

	err = store.SessionDeleteByID(session.GetID())

	if err != nil {
		t.Error("unexpected error:", err)
	}

	sessionFindWithDeleted, err := store.SessionList(SessionQuery().
		SetID(session.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(sessionFindWithDeleted) != 0 {
		t.Fatal("Session MUST be deleted, but it is not")
	}
}

func TestStore_SessionExtend(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession().
		SetValue("one two three four")

	if session == nil {
		t.Fatal("unexpected nil session")
	}

	if session.GetID() == "" {
		t.Fatal("unexpected empty id:", session.GetID())
	}

	err = store.SessionCreate(session)

	if err != nil {
		t.Error("unexpected error:", err)
	}

	if session.GetExpiresAt() == "" {
		t.Fatal("unexpected empty expiresAt:", session.GetExpiresAt())
	}

	newExpiresAt := session.GetExpiresAtCarbon().AddSeconds(100)

	if session.GetExpiresAtCarbon().Gte(newExpiresAt) {
		t.Fatal("session expiresAt must be less than new expiresAt:", newExpiresAt, " but is: ", session.GetExpiresAtCarbon())
	}

	err = store.SessionExtend(session, 100)

	if err != nil {
		t.Error("unexpected error:", err)
	}

	sessionExtended, errFind := store.SessionFindByID(session.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if sessionExtended == nil {
		t.Fatal("Session MUST NOT be nil")
	}

	if sessionExtended.GetID() != session.GetID() {
		t.Fatal("IDs do not match")
	}

	if sessionExtended.GetValue() != session.GetValue() {
		t.Fatal("Values do not match")
	}

	if sessionExtended.GetValue() != "one two three four" {
		t.Fatal("Values do not match")
	}

	if sessionExtended.GetExpiresAt() == "" {
		t.Fatal("unexpected empty expiresAt:", session.GetExpiresAt())
	}

	if sessionExtended.GetExpiresAt() == session.GetExpiresAt() {
		t.Fatal("unexpected same expiresAt:", session.GetExpiresAt())
	}

	if sessionExtended.GetExpiresAt() == sessionExtended.GetCreatedAt() {
		t.Fatal("unexpected same expiresAt:", session.GetExpiresAt())
	}

	if sessionExtended.GetExpiresAt() == sessionExtended.GetUpdatedAt() {
		t.Fatal("unexpected same expiresAt:", session.GetExpiresAt())
	}

	if sessionExtended.GetExpiresAtCarbon().Gte(newExpiresAt) {
		t.Fatal("expiresAt must be more than or equal to:", newExpiresAt, " but is: ", sessionExtended.GetExpiresAtCarbon())
	}
}

func TestStore_SessionFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession().
		SetValue("one two three four")

	if session == nil {
		t.Fatal("unexpected nil session")
	}

	if session.GetID() == "" {
		t.Fatal("unexpected empty id:", session.GetID())
	}

	err = store.SessionCreate(session)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	sessionFound, errFind := store.SessionFindByID(session.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if sessionFound == nil {
		t.Fatal("Session MUST NOT be nil")
	}

	if sessionFound.GetID() != session.GetID() {
		t.Fatal("IDs do not match")
	}

	if sessionFound.GetValue() != session.GetValue() {
		t.Fatal("Values do not match")
	}

	if sessionFound.GetValue() != "one two three four" {
		t.Fatal("Values do not match")
	}
}

func TestStore_SessionFindByKey(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession().
		SetValue("one two three four")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if session == nil {
		t.Fatal("unexpected nil session")
	}

	if session.GetKey() == "" {
		t.Fatal("unexpected empty key:", session.GetKey())
	}

	err = store.SessionCreate(session)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	sessionFound, errFind := store.SessionFindByKey(session.GetKey())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if sessionFound == nil {
		t.Fatal("Session MUST NOT be nil")
	}

	if sessionFound.GetID() != session.GetID() {
		t.Fatal("IDs do not match")
	}

	if sessionFound.GetValue() != session.GetValue() {
		t.Fatal("Values do not match")
	}

	if sessionFound.GetValue() != "one two three four" {
		t.Fatal("Values do not match")
	}
}

func TestStore_SessionList(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session1 := NewSession().
		SetUserID("1").
		SetValue("one two three")

	session2 := NewSession().
		SetUserID("2").
		SetValue("four five six")

	session3 := NewSession().
		SetUserID("3").
		SetValue("seven eight nine")

	for _, session := range []SessionInterface{session1, session2, session3} {
		err = store.SessionCreate(session)
		if err != nil {
			t.Error("unexpected error:", err)
		}
	}

	sessionList, errList := store.SessionList(SessionQuery().
		SetUserID("2").
		SetLimit(2))

	if errList != nil {
		t.Fatal("unexpected error:", errList)
	}

	if len(sessionList) != 1 {
		t.Fatal("unexpected session list length:", len(sessionList))
	}
}

func TestStore_SessionSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession()

	err = store.SessionCreate(session)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SessionSoftDeleteByID(session.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if session.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Session MUST NOT be soft deleted")
	}

	sessionFound, errFind := store.SessionFindByID(session.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if sessionFound != nil {
		t.Fatal("Session MUST be nil")
	}

	sessionFindWithSoftDeleted, err := store.SessionList(SessionQuery().
		SetID(session.GetID()).
		SetSoftDeletedIncluded(true).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(sessionFindWithSoftDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(sessionFindWithSoftDeleted[0].GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Session MUST be soft deleted", session.GetSoftDeletedAt())
	}

	if !sessionFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("Session MUST be soft deleted")
	}
}

func TestStore_SessionUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	session := NewSession()

	err = store.SessionCreate(session)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	session.SetValue("one two three")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SessionUpdate(session)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	sessionFound, errFind := store.SessionFindByID(session.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if sessionFound == nil {
		t.Fatal("Session MUST NOT be nil")
	}

	if sessionFound.GetValue() != "one two three" {
		t.Fatal("Value MUST be 'one two three', found: ", sessionFound.GetValue())
	}
}
