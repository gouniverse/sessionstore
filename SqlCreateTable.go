package sessionstore

import (
	"github.com/gouniverse/sb"
)

// SQLCreateTable returns a SQL string for creating the cache table
func (st *Store) SQLCreateTable() string {
	sql := sb.NewBuilder(st.dbDriverName).
		Table(st.sessionTableName).
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "session_key",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name:   "user_id",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 40,
		}).
		Column(sb.Column{
			Name:   "ip_address",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 50,
		}).
		Column(sb.Column{
			Name:   "user_agent",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 1024,
		}).
		Column(sb.Column{
			Name: "session_value",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: "expires_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "updated_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: "deleted_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql
}
