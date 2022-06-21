# Go MongoDB index manager

This package synchronizes MongoDB indexes with given in code.

## Usage

```
err := index_manager.SyncIndexes(context.Background(), repo.GetCollection("coll"),
    []bson.D{
        {{"token_id", int32(1)}},
        {{"token_id", int32(1)}, {"timestamp", int32(-1)}},
        {{"field", "hashed"}},
    })
if err != nil {
    panic(err.Error())
}
```
This function will create missing indexes and delete extra indexes.
Call to `SyncIndexes` is blocking.

Default "\_id\_" index is fully ignored.

It's important to use int32 values.
