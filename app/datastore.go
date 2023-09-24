package main

import (
	"sync"
	"time"
)

type RedisStore struct {
    datastore map[string]redisValue
    lock      sync.Mutex
}

type redisValue struct {
    content string
    expiration *int64
}

func NewRedisStore() *RedisStore {
    return &RedisStore{
        datastore: make(map[string]redisValue),
    }
}

func (rs *RedisStore) Set(key, value string, px *int64) {
    rs.lock.Lock()
    defer rs.lock.Unlock()
    rs.datastore[key] = newRedisValue(value, px)
}

func (rs *RedisStore) Get(key string) (string, bool) {
    rs.lock.Lock()
    defer rs.lock.Unlock()
    value, exists := rs.datastore[key]
    if value.isExpired() {
        delete(rs.datastore, key)
        return "", false
    }
    return value.content, exists
}

func newRedisValue(value string, px *int64) redisValue {
    var expirationTime int64
    if px != nil {
        currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
        expirationTime = (*px + currentTimeMillis)
    }
    return redisValue{content: value, expiration: &expirationTime}
}

func (rv redisValue) isExpired() bool {
    if rv.expiration == nil {
        return false
    }
    currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)
    return *rv.expiration < currentTimeMillis
}
