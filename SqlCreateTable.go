package sessionstore

// SQLCreateTable returns a SQL string for creating the cache table
func (st *Store) SQLCreateTable() string {
	// sql := sql.NewBuilder(st.dbDriverName).
	// 	Table(st.sessionTableName).
	// 	Column("id", sql.COLUMN_TYPE_STRING, map[string]string{
	// 		"primary": "yes",
	// 		"length":  "40",
	// 	}).
	// 	Column("session_key", sql.COLUMN_TYPE_STRING, map[string]string{}).
	// 	Column("session_value", sql.COLUMN_TYPE_TEXT, map[string]string{}).
	// 	Column("created_at", "datetime", map[string]string{}).
	// 	Column("updated_at", "datetime", map[string]string{}).
	// 	CreateIfNotExists()

	// _, err := config.Database.Exec(sql)

	// if err != nil {
	// 	cfmt.Errorln("Course table failed to be created:", err.Error())
	// }

	// return err
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.sessionTableName + ` (
	  id varchar(40) NOT NULL PRIMARY KEY,
	  session_key varchar(255) NOT NULL,
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
	  "session_key" varchar(255) NOT NULL,
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
	  "session_key" varchar(255) NOT NULL,
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
