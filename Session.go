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
		SetValue("").
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

func (o *Session) IsExpired() bool {
	return o.GetExpiresAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

func (o *Session) IsSoftDeleted() bool {
	return o.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

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
