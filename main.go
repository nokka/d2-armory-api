package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mgo "gopkg.in/mgo.v2"

	"github.com/nokka/armory/character"
	"github.com/nokka/armory/repositories"
	"github.com/nokka/armory/retrieving"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

const (
	defaultMongoDBURL = "127.0.0.1"
	defaultD2sPath    = "/Users/stekon/go/src/github.com/nokka/armory/testdata/"
	defaultDBName     = "armory"
)

func main() {
	var (
		dbname = envString("DB_NAME", defaultDBName)

		listen       = flag.String("listen", ":8090", "HTTP listen address")
		mongoDBURL   = flag.String("db.url", defaultMongoDBURL, "MongoDB URL")
		databaseName = flag.String("db.name", dbname, "MongoDB database name")
		d2spath      = flag.String("d2s.path", defaultD2sPath, "Path for parsing d2s files")
	)

	flag.Parse()

	var ctx = context.Background()

	// Logging

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewContext(logger).With("listen", *listen).With("caller", log.DefaultCaller)
	httpLogger := log.NewContext(logger).With("component", "http")

	// DB connection

	session, err := mgo.Dial(*mongoDBURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Setup repositories
	var characters character.Repository
	characters, _ = repositories.NewCharacterRepository(*databaseName, session)

	// Routing
	rs := retrieving.NewService(characters, *d2spath)

	mux := http.NewServeMux()
	mux.Handle("/retrieving/", retrieving.MakeHandler(ctx, rs, httpLogger))
	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		errs <- http.ListenAndServe(*listen, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
