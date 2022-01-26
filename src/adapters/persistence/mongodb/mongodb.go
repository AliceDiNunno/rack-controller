package mongodb

import (
	"context"
	"fmt"
	"github.com/AliceDiNunno/rack-controller/src/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func StartMongodbDatabase(config config.MongodbConfig) *mongo.Client {
	mongoURI := fmt.Sprintf("mongodb://%s:%d/", config.Host, config.Port)
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return client
}
