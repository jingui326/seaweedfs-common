package lock_manager

import (
	"github.com/seaweedfs/seaweedfs/weed/pb"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddServer(t *testing.T) {
	r := NewLockRing(100 * time.Millisecond)
	r.AddServer("localhost:8080")
	assert.Equal(t, 1, len(r.snapshots))
	r.AddServer("localhost:8081")
	r.AddServer("localhost:8082")
	r.AddServer("localhost:8083")
	r.AddServer("localhost:8084")
	r.RemoveServer("localhost:8084")
	r.RemoveServer("localhost:8082")
	r.RemoveServer("localhost:8080")

	r.RLock()
	assert.Equal(t, 8, len(r.snapshots))
	r.RUnlock()

	time.Sleep(110 * time.Millisecond)

	r.RLock()
	assert.Equal(t, 2, len(r.snapshots))
	r.RUnlock()
}

func TestLockRing(t *testing.T) {
	r := NewLockRing(100 * time.Millisecond)
	r.SetSnapshot([]pb.ServerAddress{"localhost:8080", "localhost:8081"})
	r.RLock()
	assert.Equal(t, 1, len(r.snapshots))
	r.RUnlock()

	r.SetSnapshot([]pb.ServerAddress{"localhost:8080", "localhost:8081", "localhost:8082"})
	r.RLock()
	assert.Equal(t, 2, len(r.snapshots))
	r.RUnlock()
	time.Sleep(110 * time.Millisecond)

	r.SetSnapshot([]pb.ServerAddress{"localhost:8080", "localhost:8081", "localhost:8082", "localhost:8083"})
	r.RLock()
	assert.Equal(t, 3, len(r.snapshots))
	r.RUnlock()
	time.Sleep(110 * time.Millisecond)
	r.RLock()
	assert.Equal(t, 2, len(r.snapshots))
	r.RUnlock()

	r.SetSnapshot([]pb.ServerAddress{"localhost:8080", "localhost:8081", "localhost:8082", "localhost:8083", "localhost:8084"})
	r.RLock()
	assert.Equal(t, 3, len(r.snapshots))
	r.RUnlock()
}
