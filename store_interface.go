package sessionstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)
	SessionExpiryGoroutine() error

	// New API
	SessionCount(query SessionQueryInterface) (int64, error)
	SessionCreate(session SessionInterface) error
	SessionDelete(session SessionInterface) error
	SessionDeleteByID(sessionID string) error
	SessionExtend(session SessionInterface, seconds int64) error
	SessionFindByID(sessionID string) (SessionInterface, error)
	SessionFindByKey(sessionKey string) (SessionInterface, error)
	SessionList(query SessionQueryInterface) ([]SessionInterface, error)
	SessionSoftDelete(session SessionInterface) error
	SessionSoftDeleteByID(sessionID string) error
	SessionUpdate(session SessionInterface) error
}
