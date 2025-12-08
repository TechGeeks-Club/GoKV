package store

import (
	"errors"
	"strconv"
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
	Del(keys []string) int
	Type(key string) string
	Incr(key string) int
	Exists(keys []string) int
	GetAllKeys() []string
	GetAllValues() [][]byte
}

type InMemoryStore struct {
	data map[string]KVRecord
	// TODO: add queue support
}

func NewInMemoryStore() InMemoryStore {
	return InMemoryStore{
		data: make(map[string]KVRecord),
	}
}

func (s *InMemoryStore) Set(key string, Value []byte) int {
	s.data[key] = KVRecord{Value: Value, exp: -1}
	return 1
}

func (s *InMemoryStore) Setx(key string, Value []byte, ttl int) int {
	s.data[key] = KVRecord{Value: Value, exp: int(time.Now().Unix()) + ttl}
	return 1
}

func (s *InMemoryStore) Get(key string) ([]byte, error) {
	if record, ok := s.data[key]; ok {
		return record.Value, nil
	}
	return nil, nil
}

func (s *InMemoryStore) Del(keys []string) int {
	deleted := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok {
			delete(s.data, key)
			deleted++
		}
	}
	return deleted
}
func (s *InMemoryStore) Exists(keys []string) int {
	exists := 0
	for _, key := range keys {
		if _, ok := s.data[key]; ok {
			exists++
		}
	}
	return exists
}
func (s *InMemoryStore) Incrby(key string, by int) (int, error) {
	if record, ok := s.data[key]; ok {
		rec, err := strconv.Atoi(string(record.Value))
		if err != nil {
			return 0, errors.New("ERR Not Int")
		}
		rec += by
		record.Value = []byte(strconv.Itoa(rec))
		return rec, nil
	}
	s.data[key] = KVRecord{Value: []byte(strconv.Itoa(by)), exp: -1}
	return by, nil

}

func (s *InMemoryStore) GetAllKeys() []string {
	now := int(time.Now().Unix())
	keys := make([]string, 0, len(s.data))
	for k, rec := range s.data {
		if rec.exp != -1 && rec.exp <= now {
			delete(s.data, k)
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

func (s *InMemoryStore) GetAllValues() [][]byte {
	now := int(time.Now().Unix())
	values := make([][]byte, 0, len(s.data))
	for k, rec := range s.data {
		if rec.exp != -1 && rec.exp <= now {
			delete(s.data, k)
			continue
		}
		v := make([]byte, len(rec.Value))
		copy(v, rec.Value)
		values = append(values, v)
		_ = k // keep lint happy if k unused in some builds
	}
	return values
}
