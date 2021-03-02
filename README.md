# sessionstore

# Usage

```
sessionStore = sessionstore.NewStore(sessionstore.WithGormDb(databaseInstance), sessionstore.WithTableName("milan_session"))
go sessionStore.ExpireSessionGoroutine()
	
```
