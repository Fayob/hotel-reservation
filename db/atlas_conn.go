package db

import (
	"context"
	"fmt"
	"log"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AtlasConnection() (*mongo.Client, error) {
  // Use the SetServerAPIOptions() method to set the Stable API version to 1
  serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  opts := options.Client().ApplyURI("mongodb+srv://Fabimworld:2536@atlascluster.y6qyhbt.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

  // Create a new client and connect to the server
  client, err := mongo.Connect(context.TODO(), opts)
  if err != nil {
    return nil, err
  }

  defer func() {
    if err = client.Disconnect(context.TODO()); err != nil {
      log.Fatal(err)
    }
  }()

  // Send a ping to confirm a successful connection
  // if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
  //   return nil, err
  // }
  fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client, nil
}
