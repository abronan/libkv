package etcd

import (
	"testing"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	client = "localhost:4001"
)

func makeEtcdClient(t *testing.T) store.Store {
	kv, err := New(
		[]string{client},
		&store.Config{
			ConnectionTimeout: 3 * time.Second,
			Username:          "test",
			Password:          "very-secure",
		},
	)

	if err != nil {
		t.Fatalf("cannot create store: %v", err)
	}

	return kv
}

func TestRegister(t *testing.T) {
	Register()

	kv, err := libkv.NewStore(store.ETCD, []string{client}, nil)
	assert.NoError(t, err)
	assert.NotNil(t, kv)

	if _, ok := kv.(*Etcd); !ok {
		t.Fatal("Error registering and initializing etcd")
	}
}

func TestEtcdStore(t *testing.T) {
	kv := makeEtcdClient(t)
	lockKV := makeEtcdClient(t)
	ttlKV := makeEtcdClient(t)

	testutils.RunTestCommon(t, kv)
	testutils.RunTestAtomic(t, kv)
	testutils.RunTestWatch(t, kv)
	testutils.RunTestLock(t, kv)
	testutils.RunTestLockTTL(t, kv, lockKV)
	testutils.RunTestTTL(t, kv, ttlKV)
	testutils.RunCleanup(t, kv)
}

func TestEtcdListLockKey(t *testing.T) {
	key := "testLockUnlock"
	lockKV := makeEtcdClient(t)
	value := []byte("bar")

	// We should be able to create a new lock on key
	lock, err := lockKV.NewLock(key, &store.LockOptions{Value: value, TTL: 2 * time.Second})
	assert.NoError(t, err)
	assert.NotNil(t, lock)

	pairs, err := lockKV.List(key)
	assert.NoError(t, err)
	assert.NotNil(t, pairs)
	assert.Equal(t, 1, len(pairs))
	assert.Equal(t, key, pairs[0].Key)
}
