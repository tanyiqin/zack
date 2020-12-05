package mdb

import (
	"context"
	"fmt"
	"github.com/tanyiqin/zack/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestTy1(t *testing.T) {
	p := &model.Player{
		RoleID: 123,
		Name: "qwqwqw",
	}
	data, _ := bson.Marshal(p)
	p2 := &model.Player{}
	bson.Unmarshal(data, &p2)
	fmt.Println(p2)
}

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.50.106:28017" +
	//	",192.168.50.106:28018,192.168.50.106:28019/?replicaSet=rs0"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetServerSelectionTimeout(5*time.Second))
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("test").Collection("tt")
	//result, err := coll.InsertOne(context.Background(), bson.D{
	//	{"id", 1234},
	//})
	//fmt.Println("insertOne, ",result)
	//
	//result1, _ := coll.InsertMany(context.Background(), []interface{}{
	//	bson.M{
	//		"id":444,
	//		"name":"nin",
	//	},
	//	bson.M{
	//		"id":445,
	//		"name":"wawa",
	//		"tahs":bson.A{1,2,3},
	//	},
	//})
	//fmt.Println("insertMany,", result1)

	//cursor , err := coll.Find(context.Background(), bson.D{{"id", 444}})
	//cursor.Next(ctx)
	//fmt.Println("find, ", cursor.Current.String())

	r1 :=coll.FindOne(context.Background(), bson.D{{"id", 24}})
	var temp map[string]interface{}
	err = r1.Decode(&temp)
	fmt.Println(temp, err)

}