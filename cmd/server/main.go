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
	"github.com/nokka/d2-armory-api/internal/httpserver"
	"github.com/nokka/d2-armory-api/internal/mgo"
	"github.com/nokka/d2-armory-api/internal/parsing"
	"github.com/nokka/d2-armory-api/pkg/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	var (
		httpAddress   = env.String("HTTP_ADDRESS", ":80")
		mongoDBHost   = env.String("MONGO_HOST", "mongodb:27017")
		databaseName  = env.String("MONGO_DB", "armory")
		mongoUsername = env.String("MONGO_USERNAME", "")
		mongoPassword = env.String("MONGO_PASSWORD", "")
		d2sPath       = env.String("D2S_PATH", "")
		cacheDuration = env.String("CACHE_DURATION", "3m")
	)

	if d2sPath == "" {
		log.Println("d2s path missing")
		os.Exit(0)
	}

	if mongoUsername == "" {
		log.Println("mongodb username missing")
		os.Exit(0)
	}

	if mongoPassword == "" {
		log.Println("mongodb user password missing")
		os.Exit(0)
	}

	cd, err := time.ParseDuration(cacheDuration)
	if err != nil {
		log.Printf("failed to parse cache duration, %s", err)
		os.Exit(0)
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + mongoDBHost).
		SetAuth(
			options.Credential{
				AuthSource: databaseName,
				Username:   mongoUsername,
				Password:   mongoPassword,
			})

	// Context used for mongo operations, to time them out and cancel their context.
	mgoCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err != nil {
		log.Println("failed to create context with timeout for mongodb connection", err)
		os.Exit(0)
	}

	client, err := mongo.Connect(mgoCtx, clientOptions)
	if err != nil {
		log.Println("failed to connect to mongodb", err)
		os.Exit(0)
	}

	err = client.Ping(mgoCtx, readpref.Primary())
	if err != nil {
		log.Println("failed to ping mongodb", err)
		os.Exit(0)
	}

	log.Println("connected to mongodb")

	// Repositories.
	characterRepository := mgo.NewCharacterRepository(databaseName, client)

	// Business logic services.
	parser := parsing.NewParser(d2sPath)
	characterService := character.NewService(parser, characterRepository, cd)

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
