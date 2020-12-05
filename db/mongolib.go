package mdb

import (
	"context"
	log "github.com/tanyiqin/zack/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDB struct {
	ConCancel context.CancelFunc
	Client *mongo.Client
	DB *mongo.Database
	ctx context.Context
}

var DB *MongoDB

func init() {
	var err error
	DB, err = NewMongoDB()
	if err != nil {
		log.Fatal("error connect mongodb")
	}
}

func NewMongoDB() (*MongoDB, error){
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	//Client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.50.106:28017" +
	//	",192.168.50.106:28018,192.168.50.106:28019/?replicaSet=rs0"))
	Client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017").SetServerSelectionTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	DB := Client.Database("game_server1")
	return &MongoDB{
		ConCancel: cancel,
		Client: Client,
		DB: DB,
		ctx: ctx,
	}, nil
}

func (m *MongoDB)StopMongoDB() {
	m.ConCancel()
	if err := m.Client.Disconnect(m.ctx); err != nil {
		log.Error("Client dis error", err)
	}
}

func (m *MongoDB) InsertOne(collection string, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error){
	return m.DB.Collection(collection).InsertOne(m.ctx, document, opts...)
}

func (m *MongoDB) InsertMany(collection string, documents []interface{},
	opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return m.DB.Collection(collection).InsertMany(m.ctx, documents, opts...)
}

func (m *MongoDB) FindOne(collection string, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.DB.Collection(collection).FindOne(m.ctx, filter, opts...)
}

func (m *MongoDB) Find(collection string, filter interface{},
	opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.DB.Collection(collection).Find(m.ctx, filter, opts...)
}

func (m *MongoDB) UpdateOne(collection string, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.DB.Collection(collection).UpdateOne(m.ctx, filter, update, opts...)
}

func (m *MongoDB) UpdateMany(collection string, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.DB.Collection(collection).UpdateMany(m.ctx, filter, update, opts...)
}

func (m *MongoDB) DeleteOne(collection string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.DB.Collection(collection).DeleteOne(m.ctx, filter, opts...)
}

func (m *MongoDB) DeleteMany(collection string, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.DB.Collection(collection).DeleteMany(m.ctx, filter, opts...)
}