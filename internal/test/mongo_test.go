package test

import (
	"cloud-platform-system/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestFindAndModify(t *testing.T) {
	filter := bson.D{{"admin_id", "111"}, {"user_id", "111"}, {"status", models.ApplicationFormStatusIng}}
	err := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.ApplicationFormDocument).FindOneAndUpdate(context.Background(), filter, bson.D{{"$set", bson.M{"status": 1}}}).Err()

	if err != nil && err != mongo.ErrNoDocuments {
		t.Fatal(err)
	} else if err == mongo.ErrNoDocuments {

	}
}
