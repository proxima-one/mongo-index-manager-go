# Go MongoDB index manager

This package helps you to create and delete (aka synchronize) Mongo repository 
with necessary indexes.

## Usage

```
err = index_manager.SyncIndexes(context.Background(), repo.GetCollection("coll"),
    []bson.D{
        {{"token_id", int32(1)}},
        {{"timestamp", int32(-1)}},
        {{"field", "hashed"}},
    })
if err != nil {
    panic(err.Error())
}
```
This function will create missing indexes and delete extra indexes. 
Call to `SyncIndexes` is blocking.

It's important to use int32 values.
