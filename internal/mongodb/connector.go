package mongodb

import (
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
)

// Connector ...
type Connector struct {
	session     *mgo.Session
	sessionLock sync.Mutex
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
func (c *Connector) Connect(dsn string) {
	go c.connectionLoop(dsn)
}

// GetSession ...
func (c *Connector) GetSession() *mgo.Session {
	c.sessionLock.Lock()
	defer c.sessionLock.Unlock()

	return c.session
}

// connectionLoop makes sure we're connected and informs otherwise.
func (c *Connector) connectionLoop(dsn string) {
	for {
		sess, err := mgo.Dial(dsn)
		if err == nil {
			// Connection loop that will check the state of mongoDB every 5 seconds.
			for {
				if err := sess.Ping(); err == nil {
					if c.session == nil {
						log.Println("connected to mongodb")
						c.sessionLock.Lock()
						c.session = sess
						c.sessionLock.Unlock()
					}
				} else {
					log.Println("failed to ping mongodb, removing stale session")
					sess.Close()
					c.sessionLock.Lock()
					c.session = nil
					c.sessionLock.Unlock()

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
