package shrty

import (
	"time"

	"github.com/boltdb/bolt"
)

const (
	// URLsBucket is the string name of the ShortenedURLs bucket
	URLsBucket = "ShortenedURLs"
)

// DBClient represents a client to the underlying BoltDB database
type DBClient struct {
	// Filename to the BoltDB database
	Path string

	// Returns the current time
	Now func() time.Time

	// DB itself
	db *bolt.DB
}

// NewDBClient returns a new instance of a boltDB backed client
func NewDBClient() *DBClient {
	c := &DBClient{Now: time.Now}
	return c
}

// Open opens and initializes the BoltDB database and buckets
func (c *DBClient) Open() error {
	// Open the database file
	db, err := bolt.Open(c.Path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	c.db = db

	// Initialize top-level buckets.
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte(URLsBucket)); err != nil {
		return err
	}

	return tx.Commit()
}

// Close closes the underlying BoltDB database
func (c *DBClient) Close() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}
