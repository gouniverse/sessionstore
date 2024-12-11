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

var _ SessionInterface = (*session)(nil)

// == TYPE ===================================================================

type session struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ============================================================

func NewSession() SessionInterface {
	expiresAt := carbon.Now(carbon.UTC).AddHours(2).ToDateTimeString(carbon.UTC)
	createdAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	updatedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	deletedAt := sb.MAX_DATETIME
	key := generateSessionKey(100)

	o := (&session{})

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
	o := &session{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *session) IsExpired() bool {
	return o.GetExpiresAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

func (o *session) IsSoftDeleted() bool {
	return o.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (session *session) GetCreatedAt() string {
	return session.Get(COLUMN_CREATED_AT)
}

func (session *session) GetCreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetCreatedAt(), carbon.UTC)
}

func (session *session) SetCreatedAt(createdAt string) SessionInterface {
	session.Set(COLUMN_CREATED_AT, createdAt)
	return session
}

func (session *session) GetSoftDeletedAt() string {
	return session.Get(COLUMN_SOFT_DELETED_AT)
}

func (session *session) GetSoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetSoftDeletedAt(), carbon.UTC)
}

func (session *session) SetSoftDeletedAt(DeletedAt string) SessionInterface {
	session.Set(COLUMN_SOFT_DELETED_AT, DeletedAt)
	return session
}

func (session *session) GetExpiresAt() string {
	return session.Get(COLUMN_EXPIRES_AT)
}

func (session *session) GetExpiresAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetExpiresAt(), carbon.UTC)
}

func (session *session) SetExpiresAt(expiresAt string) SessionInterface {
	session.Set(COLUMN_EXPIRES_AT, expiresAt)
	return session
}

func (session *session) GetID() string {
	return session.Get(COLUMN_ID)
}

func (session *session) SetID(id string) SessionInterface {
	session.Set(COLUMN_ID, id)
	return session
}

func (session *session) GetIPAddress() string {
	return session.Get(COLUMN_IP_ADDRESS)
}

func (session *session) SetIPAddress(iPAddress string) SessionInterface {
	session.Set(COLUMN_IP_ADDRESS, iPAddress)
	return session
}

func (session *session) GetKey() string {
	return session.Get(COLUMN_SESSION_KEY)

}

func (session *session) SetKey(key string) SessionInterface {
	session.Set(COLUMN_SESSION_KEY, key)
	return session
}

func (session *session) GetUpdatedAt() string {
	return session.Get(COLUMN_UPDATED_AT)
}

func (session *session) GetUpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetUpdatedAt(), carbon.UTC)
}

func (session *session) SetUpdatedAt(UpdatedAt string) SessionInterface {
	session.Set(COLUMN_UPDATED_AT, UpdatedAt)
	return session
}

func (session *session) GetUserAgent() string {
	return session.Get(COLUMN_USER_AGENT)
}

func (session *session) SetUserAgent(userAgent string) SessionInterface {
	session.Set(COLUMN_USER_AGENT, userAgent)
	return session
}

func (session *session) GetUserID() string {
	return session.Get(COLUMN_USER_ID)
}
func (session *session) SetUserID(userID string) SessionInterface {
	session.Set(COLUMN_USER_ID, userID)
	return session
}

func (session *session) GetValue() string {
	return session.Get(COLUMN_SESSION_VALUE)
}

func (session *session) SetValue(value string) SessionInterface {
	session.Set(COLUMN_SESSION_VALUE, value)
	return session
}
