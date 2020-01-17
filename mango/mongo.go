package mango

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//MongoConfig is config for mongodb connection
type MongoConfig struct {
	Host        string
	Port        int
	DbName      string
	User        string
	Pwd         string
	ConnTimeout int // in seconds
}

//MongoSession holds session for mongodb connection
type MongoSession struct {
	*mongo.Client
	DBname string
}

//IMongoData is interface must be implemented by all models that need to save in MongoDB
type IMongoData interface {
	GetID() interface{}
}

//InitSession initialize mongodb session
func InitSession(config MongoConfig) *MongoSession {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	return &MongoSession{client, config.DbName}
}

//Cleanup closes existing mongodb session
func (m *MongoSession) Cleanup() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	m.Disconnect(ctx)
}

//GetColl gives collection  as specified
func (m *MongoSession) GetColl(collname string) *mongo.Collection {
	return m.Client.Database(m.DBname).Collection(collname)
}

//InsertBulk inserts multiple documents at once, if document already exist, error will be raised
func (m *MongoSession) InsertBulk(col *mongo.Collection, data ...interface{}) error {
	_, err := col.InsertMany(context.Background(), data)
	if err != nil {
		return err
	}
	return nil
}

//SaveSingle insert or update single document by ID
func (m *MongoSession) SaveSingle(col *mongo.Collection, data IMongoData) error {
	filter := bson.M{"_id": data.GetID()}
	uoptn := options.Update().SetUpsert(true)
	_, err := col.UpdateOne(context.Background(), filter, data, uoptn)
	if err != nil {
		return err
	}
	return nil
}

//DeleteSingle removes single document by ID
func (m *MongoSession) DeleteSingle(col *mongo.Collection, id interface{}) error {
	filter := bson.M{"_id": id}
	_, err := col.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

// //InsertBulk inserts multiple documents at once, if document already exist, error will be raised
// func (m *MongoSession) InsertBulk(collname string, data ...interface{}) error {
// 	if data == nil {
// 		return fmt.Errorf("data is nil")
// 	}
// 	col := m.getColl(collname)

// 	_, err := col.InsertMany(context.Background(), data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// //SaveSingle insert or update single document by ID
// func (m *MongoSession) SaveSingle(collname string, data IMongoData) error {
// 	if data == nil {
// 		return fmt.Errorf("data is nil")
// 	}
// 	col := m.getColl(collname)

// 	filter := bson.M{"_id": data.GetID()}
// 	uoptn := options.Update().SetUpsert(true)
// 	_, err := col.UpdateOne(context.Background(), filter, data, uoptn)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// //DeleteSingle removes single document by ID
// func (m *MongoSession) DeleteSingle(collname string, data IMongoData) error {
// 	if data == nil {
// 		return fmt.Errorf("data is nil")
// 	}
// 	col := m.getColl(collname)

// 	filter := bson.M{"_id": data.GetID()}
// 	_, err := col.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
