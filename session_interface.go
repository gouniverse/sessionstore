package sessionstore

import "github.com/dromara/carbon/v2"

type SessionInterface interface {
	// From data object

	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Methods

	IsExpired() bool
	IsSoftDeleted() bool

	// Setters and Getters

	GetID() string
	SetID(ID string) *Session

	GetKey() string
	SetKey(Key string) *Session

	GetUserID() string
	SetUserID(UserID string) *Session

	GetIPAddress() string
	SetIPAddress(IPAddress string) *Session

	GetUserAgent() string
	SetUserAgent(UserAgent string) *Session

	GetValue() string
	SetValue(Value string) *Session

	GetExpiresAt() string
	GetExpiresAtCarbon() carbon.Carbon
	SetExpiresAt(ExpiresAt string) *Session

	GetCreatedAt() string
	GetCreatedAtCarbon() carbon.Carbon
	SetCreatedAt(createdAt string) *Session

	GetUpdatedAt() string
	GetUpdatedAtCarbon() carbon.Carbon
	SetUpdatedAt(updatedAt string) *Session

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() carbon.Carbon
	SetSoftDeletedAt(deletedAt string) *Session
}
