package storage

import (
	"sync"

	"github.com/minio/madmin-go/v2"
	"github.com/minio/minio-go/v7"
)

// Conn represents the active minio connections.
type Conn struct {

	// Client required for creating users
	admin *madmin.AdminClient

	// Client required for creating buckets
	adminClient *minio.Client

	// Ensures admin connections are thread safe.
	adminLock *sync.Mutex
}
