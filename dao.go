package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User type
// Represents a user, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type User struct {
	ID    primitive.ObjectID `bson:"_id" json:"_id"`
	Name  string             `bson:"name" json:"name"`
	Age   int                `bson:"age" json:"age"`
	Email string             `bson:"email" json:"email"`
}

// DAO declaration
type DAO struct {
	Server   string
	Database string
}

// variable declaration
var (
	db   *mongo.Database
	user User
)

// COLLECTION declaration
const (
	COLLECTION = "users"
)

// Connection to database
func (m *DAO) Connection() {
	clientOptions := options.Client().ApplyURI(m.Server)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	err = client.Connect(ctx)
	defer cancel()
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	db = client.Database(m.Database)
}

// FindAll list of users
func (m *DAO) FindAll() (users []User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "age", Value: -1}})
	cursor, err := db.Collection(COLLECTION).Find(ctx, bson.M{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
	}
	return users, err
}

// FindByID will find a user by its id
func (m *DAO) FindByID(id string) (user User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Collection(COLLECTION).FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&user)
	defer cancel()
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return user, err
}

// Delete an existing user
func (m *DAO) Delete(id string) (user User, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Collection(COLLECTION).FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&user)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	return user, err
}

// Insert a user into database
func (m *DAO) Insert(user User) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = db.Collection(COLLECTION).InsertOne(ctx, &user)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Update an existing user
func (m *DAO) Update(user User) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: user.ID}}
	update := bson.D{primitive.E{Key: "$set", Value: &user}}
	_, err = db.Collection(COLLECTION).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
