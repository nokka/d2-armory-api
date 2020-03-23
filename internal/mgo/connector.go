package mgo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connector ...
type Connector struct {
	client     *mongo.Client
	clientLock sync.Mutex
}

// NewConnector returns a new connector with all dependencies.
func NewConnector(options ...func(*Connector)) *Connector {
	c := Connector{}

	for _, op := range options {
		op(&c)
	}

	return &c
}

// Connect ...
func (c *Connector) Connect(ctx context.Context, dsn string) {
	go c.connectionLoop(ctx, dsn)
}

// GetClient ...
func (c *Connector) GetClient() *mongo.Client {
	c.clientLock.Lock()
	defer c.clientLock.Unlock()

	return c.client
}

// connectionLoop makes sure we're connected and informs otherwise.
func (c *Connector) connectionLoop(ctx context.Context, dsn string) {
	for {
		fmt.Println("TRYING TO CONNECT")
		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://armory:not_secure_at_all@localhost:27017"))
		//client, err := mongo.NewClient(options.Client().ApplyURI(dsn))
		if err != nil {
			log.Println("failed to create mongo client", err)
		}

		err = client.Connect(ctx)
		if err != nil {
			log.Println("failed to connect mongo client", err)
		}

		if err == nil {
			// Connection loop that will check the state of mongoDB every 5 seconds.
			for {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				if err = client.Ping(ctx, readpref.Primary()); err == nil {
					if c.client == nil {
						log.Println("connected to mongodb")
						c.clientLock.Lock()
						c.client = client
						c.clientLock.Unlock()
					}
				} else {
					log.Println("failed to ping mongodb:", err)
					client.Disconnect(ctx)
					c.clientLock.Lock()
					c.client = nil
					c.clientLock.Unlock()

					break
				}

				time.Sleep(5 * time.Second)
			}
		} else {
			log.Printf("failed to connect to mongodb: %s", err)
		}

		time.Sleep(5 * time.Second)
	}
}
