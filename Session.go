package sessionstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

type SessionOptions struct {
	UserID    string
	IPAddress string
	UserAgent string
}

var _ SessionInterface = (*Session)(nil)

// == TYPE ===================================================================

type Session struct {
	// id        string     `db:"id"`            // varchar(40), primary key
	// key       string     `db:"session_key"`   // varchar(40)
	// userID    string     `db:"user_id"`       // varchar(40)
	// iPAddress string     `db:"ip_address"`    // varchar(50)
	// userAgent string     `db:"user_agent"`    // varchar(1024)
	// value     string     `db:"session_value"` // long text
	// expiresAt *time.Time `db:"expires_at"`    // datetime NOT NULL
	// createdAt time.Time  `db:"created_at"`    // datetime NOT NULL
	// updatedAt time.Time  `db:"updated_at"`    // datetime NOT NULL
	// deletedAt *time.Time `db:"deleted_at"`    // datetime DEFAULT NULL
	dataobject.DataObject
}

// == CONSTRUCTORS ============================================================

func NewSession() SessionInterface {
	expiresAt := carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC)
	createdAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	updatedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	deletedAt := sb.MAX_DATETIME
	key := generateSessionKey(100)

	o := (&Session{})

	o.SetID(uid.HumanUid()).
		SetKey(key).
		SetUserID("").
		SetUserAgent("").
		SetIPAddress("").
		SetExpiresAt(expiresAt).
		SetCreatedAt(createdAt).
		SetUpdatedAt(updatedAt).
		SetSoftDeletedAt(deletedAt)

	return o
}

func NewSessionFromExistingData(data map[string]string) SessionInterface {
	o := &Session{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *Session) IsSoftDeleted() bool {
	return o.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
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

// == SETTERS AND GETTERS =====================================================

func (session *Session) GetCreatedAt() string {
	return session.Get(COLUMN_CREATED_AT)
}

func (session *Session) GetCreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetCreatedAt(), carbon.UTC)
}

func (session *Session) SetCreatedAt(createdAt string) *Session {
	session.Set(COLUMN_CREATED_AT, createdAt)
	return session
}

func (session *Session) GetSoftDeletedAt() string {
	return session.Get(COLUMN_SOFT_DELETED_AT)
}

func (session *Session) GetSoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetSoftDeletedAt(), carbon.UTC)
}

func (session *Session) SetSoftDeletedAt(DeletedAt string) *Session {
	session.Set(COLUMN_SOFT_DELETED_AT, DeletedAt)
	return session
}

func (session *Session) GetExpiresAt() string {
	return session.Get(COLUMN_EXPIRES_AT)
}

func (session *Session) GetExpiresAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetExpiresAt(), carbon.UTC)
}

func (session *Session) SetExpiresAt(expiresAt string) *Session {
	session.Set(COLUMN_EXPIRES_AT, expiresAt)
	return session
}

func (session *Session) GetID() string {
	return session.Get(COLUMN_ID)
}

func (session *Session) SetID(id string) *Session {
	session.Set(COLUMN_ID, id)
	return session
}

func (session *Session) GetIPAddress() string {
	return session.Get(COLUMN_IP_ADDRESS)
}

func (session *Session) SetIPAddress(iPAddress string) *Session {
	session.Set(COLUMN_IP_ADDRESS, iPAddress)
	return session
}

func (session *Session) GetKey() string {
	return session.Get(COLUMN_SESSION_KEY)

}

func (session *Session) SetKey(key string) *Session {
	session.Set(COLUMN_SESSION_KEY, key)
	return session
}

func (session *Session) GetUpdatedAt() string {
	return session.Get(COLUMN_UPDATED_AT)
}

func (session *Session) GetUpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetUpdatedAt(), carbon.UTC)
}

func (session *Session) SetUpdatedAt(UpdatedAt string) *Session {
	session.Set(COLUMN_UPDATED_AT, UpdatedAt)
	return session
}

func (session *Session) GetUserAgent() string {
	return session.Get(COLUMN_USER_AGENT)
}

func (session *Session) SetUserAgent(userAgent string) *Session {
	session.Set(COLUMN_USER_AGENT, userAgent)
	return session
}

func (session *Session) GetUserID() string {
	return session.Get(COLUMN_USER_ID)
}
func (session *Session) SetUserID(userID string) *Session {
	session.Set(COLUMN_USER_ID, userID)
	return session
}
func (session *Session) GetValue() string {
	return session.Get(COLUMN_SESSION_VALUE)
}

func (session *Session) SetValue(value string) *Session {
	session.Set(COLUMN_SESSION_VALUE, value)
	return session
}
