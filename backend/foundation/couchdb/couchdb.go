// Package couchdb provides support for access the database.
package couchdb

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-kivik/couchdb/v4" // The CouchDB driver
	kivik "github.com/go-kivik/kivik/v4"
)

// Config is the required properties to use the database.
type Config struct {
	Protocol	string
	User 		string
	Password	string
	Host		string
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*kivik.Client, error) {
	fmt.Println(cfg.Protocol + "://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host)
	client, err := kivik.New("couch",
		cfg.Protocol + "://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host)
	if err != nil {
		return nil, err
	}
	fmt.Println("connected to database")
	return client, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *kivik.Client) error {

	// First check we can ping the database.
	for attempts := 1; ; attempts++ {
		up, pingError := db.Ping(ctx)
		if up && pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	// Make sure we didn't timeout or be cancelled.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	return nil
}

