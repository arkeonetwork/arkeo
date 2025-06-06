package sentinel

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNonceStore_InMemory(t *testing.T) {
	// Create in-memory store
	store, err := NewNonceStore("")
	require.NoError(t, err)
	defer store.Close()

	// Test Get on non-existent key
	nonce, err := store.Get(12345)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), nonce)

	// Test Set and Get
	err = store.Set(12345, 42)
	assert.NoError(t, err)

	nonce, err = store.Get(12345)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), nonce)

	// Test Update
	err = store.Set(12345, 100)
	assert.NoError(t, err)

	nonce, err = store.Get(12345)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), nonce)

	// Test multiple contracts
	err = store.Set(67890, 1)
	assert.NoError(t, err)

	nonce1, err := store.Get(12345)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), nonce1)

	nonce2, err := store.Get(67890)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), nonce2)
}

func TestNonceStore_Persistent(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "nonce-store-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create store and set value
	store1, err := NewNonceStore(tmpDir)
	require.NoError(t, err)

	err = store1.Set(12345, 42)
	assert.NoError(t, err)

	// Close store
	err = store1.Close()
	assert.NoError(t, err)

	// Reopen store and check value persisted
	store2, err := NewNonceStore(tmpDir)
	require.NoError(t, err)
	defer store2.Close()

	nonce, err := store2.Get(12345)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), nonce)
}

func TestNonceStore_ConcurrentAccess(t *testing.T) {
	store, err := NewNonceStore("")
	require.NoError(t, err)
	defer store.Close()

	// Concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			err := store.Set(uint64(i), int64(i*10))
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all writes
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values
	for i := 0; i < 10; i++ {
		nonce, err := store.Get(uint64(i))
		assert.NoError(t, err)
		assert.Equal(t, int64(i*10), nonce)
	}
}

func TestNonceRecord_UpdatedAt(t *testing.T) {
	store, err := NewNonceStore("")
	require.NoError(t, err)
	defer store.Close()

	before := time.Now().Unix()
	err = store.Set(12345, 42)
	assert.NoError(t, err)
	after := time.Now().Unix()

	// Get raw record to check UpdatedAt
	key := "12345"
	value, err := store.db.Get([]byte(key), nil)
	assert.NoError(t, err)

	var record NonceRecord
	err = json.Unmarshal(value, &record)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, record.UpdatedAt, before)
	assert.LessOrEqual(t, record.UpdatedAt, after)
}