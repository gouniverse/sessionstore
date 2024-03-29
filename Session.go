package sessionstore

import (
	"time"
)

// Session type
type Session struct {
	ID        string     `db:"id"`            // varchar(40), primary key
	Key       string     `db:"session_key"`   // varchar(40)
	UserID    string     `db:"user_id"`       // varchar(40)
	IPAddress string     `db:"ip_address"`    // varchar(50)
	UserAgent string     `db:"user_agent"`    // varchar(1024)
	Value     string     `db:"session_value"` // long text
	ExpiresAt *time.Time `db:"expires_at"`    // datetime NOT NULL
	CreatedAt time.Time  `db:"created_at"`    // datetime NOT NULL
	UpdatedAt time.Time  `db:"updated_at"`    // datetime NOT NULL
	DeletedAt *time.Time `db:"deleted_at"`    // datetime DEFAULT NULL
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

// // SessionFindByToken finds a session by key
// func SessionFindByToken(key string) *Session {
// 	session := &Session{}
// 	result := GetDb().Where("`key` = ?", key).First(&session)

// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return nil
// 	}

// 	return session
// }

// // SessionGet gets a key from cache
// func SessionGet(key string, valueDefault string) string {
// 	cache := SessionFindByToken(key)

// 	if cache != nil {
// 		return cache.Value
// 	}

// 	return valueDefault
// }

// // SessionStart starts a session with a specified key
// func SessionStart(key string) bool {
// 	session := SessionFindByToken(key)
// 	expiresAt := time.Now().Add(time.Duration(SessionTimeoutSeconds) * time.Second)

// 	if session != nil {
// 		return true
// 	}

// 	var newSession = Session{Key: key, Value: "{}", ExpiresAt: &expiresAt}

// 	dbResult := GetDb().Create(&newSession)

// 	if dbResult.Error != nil {
// 		return false
// 	}

// 	return true
// }

// // SessionSet sets a key in cache
// func SessionSet(key string, value string, seconds int64) bool {
// 	session := SessionFindByToken(key)
// 	expiresAt := time.Now().Add(time.Duration(SessionTimeoutSeconds) * time.Duration(seconds))

// 	if session != nil {
// 		session.Value = value
// 		session.ExpiresAt = &expiresAt
// 		//dbResult := GetDb().Table(User).Where("`key` = ?", key).Update(&cache)
// 		dbResult := GetDb().Save(&session)
// 		if dbResult != nil {
// 			return false
// 		}
// 		return true
// 	}

// 	var newSessiom = Session{Key: key, Value: value, ExpiresAt: &expiresAt}

// 	dbResult := GetDb().Create(&newSessiom)

// 	if dbResult.Error != nil {
// 		return false
// 	}

// 	return true
// }

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

// // SessionEmpty removes all keys from the sessiom
// func SessionEmpty(sessionKey string) bool {
// 	session := SessionFindByToken(sessionKey)

// 	kv := hashmap.New()

// 	if session == nil {
// 		return true
// 	}

// 	json, err := kv.ToJSON()

// 	if err != nil {
// 		return false
// 	}

// 	session.Value = string(json)

// 	GetDb().Save(&session)

// 	return true
// }

// // SessionSetKey gets a key from sessiom
// func SessionSetKey(sessionKey string, key string, value string) bool {
// 	session := SessionFindByToken(sessionKey)

// 	kv := hashmap.New()

// 	if session == nil {
// 		isOk := SessionSet(sessionKey, "{}", 2000)
// 		if isOk == false {
// 			return false
// 		}
// 		session = SessionFindByToken(sessionKey)
// 	}

// 	log.Println(value)

// 	err := kv.FromJSON([]byte(session.Value))
// 	if err != nil {
// 		return false
// 	}

// 	kv.Put(key, value)
// 	json, err := kv.ToJSON()

// 	if err != nil {
// 		return false
// 	}
// 	log.Println(string(json))

// 	session.Value = string(json)

// 	GetDb().Save(&session)

// 	return true
// }

// // SessionRemoveKey removes a key from sessiom
// func SessionRemoveKey(sessionKey string, key string) bool {
// 	session := SessionFindByToken(sessionKey)

// 	kv := hashmap.New()

// 	if session == nil {
// 		return true
// 	}

// 	err := kv.FromJSON([]byte(session.Value))
// 	if err != nil {
// 		return false
// 	}

// 	kv.Remove(key)

// 	json, err := kv.ToJSON()

// 	if err != nil {
// 		return false
// 	}

// 	log.Println(string(json))

// 	session.Value = string(json)

// 	GetDb().Save(&session)

// 	return true
// }

// // SessionExpireJobGoroutine - soft deletes expired cache
// func SessionExpireJobGoroutine() {
// 	i := 0
// 	for {
// 		i++
// 		log.Println("Cleaning expired sessions...")
// 		GetDb().Where("`expires_at` < ?", time.Now()).Delete(Session{})
// 		time.Sleep(60 * time.Second) // Every minute
// 	}
// }
