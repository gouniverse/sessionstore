# sessionstore

# Usage

```
sessionStore = sessionstore.NewStore(sessionstore.WithGormDb(databaseInstance), sessionstore.WithTableName("my_session"))

go sessionStore.ExpireSessionGoroutine()
```
