package viewmodels

import (
	"context"
	"fmt"
	"log"

	"github.com/samtech09/dbtools/mango"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

//DbUser is sample database user
type DbUser struct {
	ID   int
	Name string
}

type dbUserCache struct {
	m        *mango.MongoSession
	c        *mongo.Collection
	collname string
}

//GetID implements interface to qualify for MongoDb cache
func (d DbUser) GetID() interface{} {
	return d.ID
}

func (c *dbUserCache) toInterface(list []DbUser) []interface{} {
	var islice []interface{} = make([]interface{}, len(list))
	for i, v := range list {
		islice[i] = v
	}
	return islice
}

//Cache returns mongodb collection for Dbuser
func (d *DbUser) Cache(m *mango.MongoSession) *dbUserCache {
	colname := "dbusers"
	cac := dbUserCache{m, m.GetColl(colname), colname}
	return &cac
}

//GetAll returns all documents from given collection
func (c *dbUserCache) GetAll() ([]*DbUser, error) {
	cur, err := c.c.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	// Close the cursor once finished
	defer cur.Close(context.TODO())

	results := make([]*DbUser, 0, 10)
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem DbUser
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

//SaveDbUsers insert or update one or more MongoData to database
func (c *dbUserCache) InsertBulk(data ...DbUser) error {
	if data == nil {
		return fmt.Errorf("nothing to save")
	}
	is := c.toInterface(data)
	return c.m.InsertBulk(c.c, is)
}

//Upsert insert or update single document by ID
func (c *dbUserCache) Upsert(data DbUser) error {
	c.m.SaveSingle(c.c, data)
	return nil
}

//Upsert insert or update single document by ID
func (c *dbUserCache) Delete(id int) error {
	c.m.DeleteSingle(c.c, id)
	return nil
}
