package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nokka/d2-armory-api/internal/character"
	"github.com/nokka/d2-armory-api/internal/domain"
	"github.com/nokka/d2-armory-api/internal/httpserver"
	"github.com/nokka/d2-armory-api/internal/mongodb"
	"github.com/nokka/d2-armory-api/pkg/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	var (
		httpAddress = env.String("HTTP_ADDRESS", ":80")
		//mongoDBHost   = env.String("MONGO_HOST", "mongodb:27017")
		databaseName  = env.String("MONGO_DB", "armory")
		mongoUsername = env.String("MONGO_USERNAME", "")
		mongoPassword = env.String("MONGO_PASSWORD", "")
		d2sPath       = env.String("D2S_PATH", "")
		cacheDuration = env.String("CACHE_DURATION", "3m")
	)

	if d2sPath == "" {
		log.Println("D2S path missing")
		os.Exit(0)
	}

	if mongoUsername == "" {
		log.Println("Mongodb username missing")
		os.Exit(0)
	}

	if mongoPassword == "" {
		log.Println("Mongodb user password missing")
		os.Exit(0)
	}

	cd, err := time.ParseDuration(cacheDuration)
	if err != nil {
		log.Printf("failed to parse cache duration, %s", err)
		os.Exit(0)
	}

	// Setup MongoDB.
	c := mongodb.NewConnector()

	//dsn := fmt.Sprintf("mongodb://%s:%s@%s", mongoUsername, mongoPassword, mongoDBHost)

	/*c.Connect(fmt.Sprintf(
		"mongodb://%s:%s@%s/%s",
		mongoUsername,
		mongoPassword,
		mongoDBHost,
		databaseName,
	))*/
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").
		SetAuth(
			options.Credential{
				AuthSource: "armory",
				Username:   "armory",
				Password:   "not_secure_at_all",
			})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}

	go func() {
		for {
			fmt.Println("starting ping")
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			err := client.Ping(ctx, readpref.Primary())
			fmt.Println("ping error", err)

			if err == nil {
				var result domain.Character
				filter := bson.M{"id": "cathans"}
				findCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				collection := client.Database("armory").Collection("character")
				err = collection.FindOne(findCtx, filter).Decode(&result)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(result)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// NEW Mongodb.
	//cl := mgo.NewConnector()
	//cl.Connect(context.Background(), dsn)

	// Repositories.
	characterRepository := mongodb.NewCharacterRepository(databaseName, c)

	// Services.
	characterService := character.NewService(d2sPath, characterRepository, cd)

	// Channel to receive errors on.
	errorChannel := make(chan error)

	// HTTP server.
	go func() {
		httpServer := httpserver.NewServer(httpAddress, characterService)
		errorChannel <- httpServer.Open()
	}()

	// Capture interupts.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errorChannel <- fmt.Errorf("got signal %s", <-c)
	}()

	// Listen for errors indefinitely.
	if err := <-errorChannel; err != nil {
		os.Exit(1)
	}
}
