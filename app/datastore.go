package main

import "sync"

type RedisStore struct {
    datastore map[string]string
    lock      sync.Mutex
}

func NewRedisStore() *RedisStore {
    return &RedisStore{
        datastore: make(map[string]string),
    }
}

func (rs *RedisStore) Set(key, value string) {
    rs.lock.Lock()
    defer rs.lock.Unlock()
    rs.datastore[key] = value
}

func (rs *RedisStore) Get(key string) (string, bool) {
    rs.lock.Lock()
    defer rs.lock.Unlock()
    value, exists := rs.datastore[key]
    return value, exists
}
