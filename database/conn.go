package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db = "mongodb+srv://durgeshchaudhary020401:Durgesh%402022@cluster0.vzdrcpb.mongodb.net/durgesh"
var SaveData *mongo.Collection

func ConnectDB() error {

	clientOption := options.Client().ApplyURI(db)
	client, err := mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		return err
	}
	fmt.Println("Database connected")
	SaveData = client.Database("Player").Collection("Save")

	return nil
}
