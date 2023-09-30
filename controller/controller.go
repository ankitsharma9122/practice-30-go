package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

const mongoURI = "mongodb+srv://mongo_user:2le4EJFkGBN35HYN@cluster0.wxybdpj.mongodb.net/mydbwriters?retryWrites=true&w=majority"
const dbName = ""

func init() {
	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
	Client, _ = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
}

func ConnectionDbInstance() *mongo.Collection {
	collection := Client.Database("mydbwriters").Collection("blogswriter")
	return collection
}
