package sessionstore

type StoreInterface interface {
	EnableDebug(debug bool)
	AutoMigrate() error

	// Old API
	Extend(sessionKey string, seconds int64, options SessionOptions) error
	Delete(sessionKey string, options SessionOptions) error
	Get(key string, valueDefault string, options SessionOptions) (string, error)
	GetAny(key string, valueDefault interface{}, options SessionOptions) (interface{}, error)
	GetMap(key string, valueDefault map[string]any, options SessionOptions) (map[string]any, error)
	MergeMap(key string, mergeMap map[string]any, seconds int64, options SessionOptions) error
	Set(key string, value string, seconds int64, options SessionOptions) error
	SetAny(key string, value interface{}, seconds int64, options SessionOptions) error
	SetMap(key string, value map[string]any, seconds int64, options SessionOptions) error

	// New API
	SessionCount(query SessionQueryInterface) (int64, error)
	SessionCreate(session SessionInterface) error
	SessionList(query SessionQueryInterface) ([]SessionInterface, error)
	SessionUpdate(session SessionInterface, options SessionOptions) error
	// SessionDelete(session Session, options SessionOptions) error
}
