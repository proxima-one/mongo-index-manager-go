package index_manager

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func isEqualE(a, b bson.E) bool {
	return a.Key == b.Key && a.Value == b.Value
}

func isEqual(a, b bson.D) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !isEqualE(a[i], b[i]) {
			return false
		}
	}
	return true
}

func contains(indexes []bson.D, index bson.D) bool {
	for _, ind := range indexes {
		if isEqual(ind, index) {
			return true
		}
	}
	return false
}

func SyncIndexes(ctx context.Context, collection *mongo.Collection, requiredIndexes []bson.D) error {
	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return err
	}
	var existingIndexesStructs []struct {
		Name string
		Key  bson.D
	}
	err = cur.All(ctx, &existingIndexesStructs)
	if err != nil {
		return err
	}
	var existingIndexes []bson.D
	for _, index := range existingIndexesStructs {
		if index.Name == "_id_" {
			continue
		}
		existingIndexes = append(existingIndexes, index.Key)
	}

	var indexesToDelete []string
	for _, existingIndex := range existingIndexesStructs {
		if existingIndex.Name == "_id_" {
			continue
		}
		if !contains(requiredIndexes, existingIndex.Key) {
			indexesToDelete = append(indexesToDelete, existingIndex.Name)
		}
	}

	var indexesToCreate []mongo.IndexModel
	for _, requiredIndex := range requiredIndexes {
		if !contains(existingIndexes, requiredIndex) {
			indexesToCreate = append(indexesToCreate, mongo.IndexModel{Keys: requiredIndex})
		}
	}

	if len(indexesToCreate) > 0 {
		_, err = collection.Indexes().CreateMany(ctx, indexesToCreate)
		if err != nil {
			return err
		}
	}

	for _, index := range indexesToDelete {
		_, err = collection.Indexes().DropOne(ctx, index)
		if err != nil {
			return err
		}
	}
	return nil
}
