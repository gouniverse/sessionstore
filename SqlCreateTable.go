package sessionstore

// SQLCreateTable returns a SQL string for creating the cache table
func (st *Store) SQLCreateTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.sessionTableName + ` (
	  id varchar(40) NOT NULL PRIMARY KEY,
	  session_key varchar(40) NOT NULL,
	  session_value text,
	  expires_at datetime,
	  created_at datetime NOT NULL,
	  updated_at datetime NOT NULL,
	  deleted_at datetime
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.sessionTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "session_key" varchar(40) NOT NULL,
	  "session_value" text,
	  "expires_at" timestamptz(6),
	  "created_at" timestamptz(6) NOT NULL,
	  "updated_at" timestamptz(6) NOT NULL,
	  "deleted_at" timestamptz(6)
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.sessionTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "session_key" varchar(40) NOT NULL,
	  "session_value" text,
	  "expires_at" datetime,
	  "created_at" datetime NOT NULL,
	  "updated_at" datetime NOT NULL,
	  "deleted_at" datetime
	)
	`

	sql := "unsupported driver " + st.dbDriverName

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}
