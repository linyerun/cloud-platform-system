package test

import (
	"cloud-platform-system/internal/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestFindAndModify(t *testing.T) {
	filter := bson.D{{"admin_id", "111"}, {"user_id", "111"}, {"status", models.UserApplicationFormStatusIng}}
	err := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.UserApplicationFormDocument).FindOneAndUpdate(context.Background(), filter, bson.D{{"$set", bson.M{"status": 1}}}).Err()

	if err != nil && err != mongo.ErrNoDocuments {
		t.Fatal(err)
	} else if err == mongo.ErrNoDocuments {

	}
}

func TestFindTask(t *testing.T) {
	filter := bson.M{"_id": "KaTAI1XAAAA="}
	result := svcCtx.MongoClient.Database(svcCtx.Config.Mongo.DbName).Collection(models.AsyncTaskDocument).FindOne(context.Background(), filter)
	if err := result.Err(); err != nil && err != mongo.ErrNoDocuments {
		t.Fatal(err)
	} else if err == mongo.ErrNoDocuments {
		fmt.Println("empty")
		return
	}
	asyncTask := new(models.AsyncTask)
	err := result.Decode(asyncTask)
	if err != nil {
		t.Fatal(err)
	}
}
