package session

import (
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

// Session type
type Session struct {
	ID        string     `gorm:"type:varchar(40);column:id;primary_key;"`
	Key       string     `gorm:"type:varchar(40);column:key;"`
	Value     string     `gorm:"type:longtext;column:value;"`
	ExpiresAt *time.Time `gorm:"type:datetime;olumn:expores_at;DEFAULT NULL;"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt *time.Time `gorm:"type:datetime;olumn:deleted_at;DEFAULT NULL;"`
}

// BeforeCreate adds UID to model
func (c *Session) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.NanoUid()
	c.ID = uuid
	return nil
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
