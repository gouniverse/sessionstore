package sessionstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

var _ SessionInterface = (*session)(nil)

// == TYPE ===================================================================

// session represents a user session.
type session struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ============================================================

// NewSession creates a new session.
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

// NewSessionFromExistingData creates a new session from existing data.
func NewSessionFromExistingData(data map[string]string) SessionInterface {
	o := &session{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

// IsExpired returns true if the session is expired
func (o *session) IsExpired() bool {
	return o.GetExpiresAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// IsSoftDeleted returns true if the session is soft deleted
func (o *session) IsSoftDeleted() bool {
	return o.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

// GetCreatedAt returns the created at time of the session
func (session *session) GetCreatedAt() string {
	return session.Get(COLUMN_CREATED_AT)
}

// GetCreatedAtCarbon returns the created at time of the session as a carbon object
func (session *session) GetCreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetCreatedAt(), carbon.UTC)
}

// SetCreatedAt sets the created at time of the session
func (session *session) SetCreatedAt(createdAt string) SessionInterface {
	session.Set(COLUMN_CREATED_AT, createdAt)
	return session
}

// GetSoftDeletedAt returns the soft deleted at time of the session
func (session *session) GetSoftDeletedAt() string {
	return session.Get(COLUMN_SOFT_DELETED_AT)
}

// GetSoftDeletedAtCarbon returns the soft deleted at time of the session as a carbon object.
func (session *session) GetSoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetSoftDeletedAt(), carbon.UTC)
}

// SetSoftDeletedAt sets the soft deleted at time of the session.
func (session *session) SetSoftDeletedAt(DeletedAt string) SessionInterface {
	session.Set(COLUMN_SOFT_DELETED_AT, DeletedAt)
	return session
}

// GetExpiresAt returns the expires at time of the session.
func (session *session) GetExpiresAt() string {
	return session.Get(COLUMN_EXPIRES_AT)
}

// GetExpiresAtCarbon returns the expires at time of the session as a carbon object.
func (session *session) GetExpiresAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetExpiresAt(), carbon.UTC)
}

// SetExpiresAt sets the expires at time of the session.
func (session *session) SetExpiresAt(expiresAt string) SessionInterface {
	session.Set(COLUMN_EXPIRES_AT, expiresAt)
	return session
}

// GetID returns the id of the session.
func (session *session) GetID() string {
	return session.Get(COLUMN_ID)
}

// SetID sets the id of the session.
func (session *session) SetID(id string) SessionInterface {
	session.Set(COLUMN_ID, id)
	return session
}

// GetIPAddress returns the IP address of the session.
func (session *session) GetIPAddress() string {
	return session.Get(COLUMN_IP_ADDRESS)
}

// SetIPAddress sets the IP address of the session.
func (session *session) SetIPAddress(iPAddress string) SessionInterface {
	session.Set(COLUMN_IP_ADDRESS, iPAddress)
	return session
}

// GetKey returns the key of the session.
func (session *session) GetKey() string {
	return session.Get(COLUMN_SESSION_KEY)

}

// SetKey sets the key of the session.
func (session *session) SetKey(key string) SessionInterface {
	session.Set(COLUMN_SESSION_KEY, key)
	return session
}

// GetUpdatedAt returns the updated at time of the session.
func (session *session) GetUpdatedAt() string {
	return session.Get(COLUMN_UPDATED_AT)
}

// GetUpdatedAtCarbon returns the updated at time of the session as a carbon object.
func (session *session) GetUpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(session.GetUpdatedAt(), carbon.UTC)
}

// SetUpdatedAt sets the updated at time of the session.
func (session *session) SetUpdatedAt(UpdatedAt string) SessionInterface {
	session.Set(COLUMN_UPDATED_AT, UpdatedAt)
	return session
}

// GetUserAgent returns the user agent of the session.
func (session *session) GetUserAgent() string {
	return session.Get(COLUMN_USER_AGENT)
}

// SetUserAgent sets the user agent of the session.
func (session *session) SetUserAgent(userAgent string) SessionInterface {
	session.Set(COLUMN_USER_AGENT, userAgent)
	return session
}

// GetUserID returns the user id of the session.
func (session *session) GetUserID() string {
	return session.Get(COLUMN_USER_ID)
}

// SetUserID sets the user id of the session.
func (session *session) SetUserID(userID string) SessionInterface {
	session.Set(COLUMN_USER_ID, userID)
	return session
}

// GetValue returns the value of the session.
func (session *session) GetValue() string {
	return session.Get(COLUMN_SESSION_VALUE)
}

// SetValue sets the value of the session.
func (session *session) SetValue(value string) SessionInterface {
	session.Set(COLUMN_SESSION_VALUE, value)
	return session
}
