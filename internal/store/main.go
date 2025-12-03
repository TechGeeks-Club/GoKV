package store

import (
	"time"
)

// for this version we will just keep it simple (map)

type KVRecord struct {
	Value []byte
	exp   int // exp = unix_time_now + ttl
}

type KVStore interface {
	Set(Key string, Value []byte) int
	Setx(Key string, Value []byte, exp int) int
	Get(key string) ([]byte, error)
	Delete(key string) int
}

type InMemoryStore struct {
	data map[string]KVRecord
	// TODO: add queue support
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]KVRecord),
	}
}

func (s *InMemoryStore) Set(key string, Value []byte) int {
	s.data[key] = KVRecord{Value: Value, exp: -1}
	return 1
}

func (s *InMemoryStore) Setx(key string, Value []byte, ttl int) int {
	now := time.Now()
	unixTimeSeconds := now.Unix()
	s.data[key] = KVRecord{Value: Value, exp: unixTimeSeconds + ttl}
	return 1
}

func (s *InMemoryStore) Get(key string) ([]byte, error) {
	if record, ok := s.data[key]; ok {
		return record.Value, nil
	}
	return nil, nil
}

func (s *InMemoryStore) Delete(key string) int {
	if _, ok := s.data[key]; ok {
		delete(s.data, key)
		return 1
	}
	return 0
}
