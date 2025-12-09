package store

import (
	"strconv"
	"time"

	"github.com/B-AJ-Amar/gokv/internal/common"
)

// for this version we will just keep it simple (map)

const (
	ExpireNone = iota
	ExpireEX
	ExpirePX
	ExpireEXAT
	ExpirePXAT
	ExpireKEEPTTL
)

type SetArgs struct {
	ExpType int8
	ExpVal  int
	NX_XX   int8 // 0 for nil, 1 for NX, 2 for XX
	KeepTTL bool
	Get     bool
}

type KVRecord struct {
	Value []byte
	exp   int64 // exp = unix_time_now(ms) + ttl(ms)
}

type KVStore interface {
	Set(Key string, Value []byte) int
	Setx(key string, Value []byte, args SetArgs) (int, []byte, error)
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

// TODO: add mutex for set and setx
func (s *InMemoryStore) Set(key string, Value []byte) int {
	s.data[key] = KVRecord{Value: Value, exp: -1}
	return 1
}

func (s *InMemoryStore) Setx(key string, Value []byte, args SetArgs) (int, []byte, error) {
	expUnix := int64(-1)
	oldValue := []byte{}
	retOld := false

	switch args.NX_XX {
	case 1: // NX
		if _, ok := s.data[key]; ok {
			return 0, nil, nil
		}
		if args.Get {
			return 0, nil, common.ErrSyntaxError
		}
	case 2: // XX
		if _, ok := s.data[key]; !ok {
			return 0, nil, nil
		}
	}

	nowMs := time.Now().UnixMilli()
	switch args.ExpType {
	case ExpireEX:
		expUnix = nowMs + int64(args.ExpVal)*1000 // EX is seconds, convert to ms
	case ExpirePX:
		expUnix = nowMs + int64(args.ExpVal) // PX is ms
	case ExpireEXAT:
		expUnix = int64(args.ExpVal) * 1000 // EXAT is seconds, convert to ms
	case ExpirePXAT:
		expUnix = int64(args.ExpVal) // PXAT is ms
	default:
		expUnix = -1
	}

	if args.KeepTTL {
		if record, ok := s.data[key]; ok {
			expUnix = record.exp
		}
	}

	if args.Get {
		if record, ok := s.data[key]; ok {
			oldValue = record.Value
			retOld = true
		}
	}

	s.data[key] = KVRecord{Value: Value, exp: expUnix}
	if retOld {
		return 1, oldValue, nil
	}
	return 1, nil, nil
}

func (s *InMemoryStore) Get(key string) ([]byte, error) {
	if record, ok := s.data[key]; ok {
		nowMs := time.Now().UnixMilli()
		if record.exp != -1 && record.exp <= nowMs {
			delete(s.data, key)
			return nil, nil
		}
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
			return 0, common.ErrNotIntOROutOfRange
		}
		rec += by
		record.Value = []byte(strconv.Itoa(rec))
		return rec, nil
	}
	s.data[key] = KVRecord{Value: []byte(strconv.Itoa(by)), exp: -1}
	return by, nil

}

func (s *InMemoryStore) GetAllKeys() []string {
	nowMs := time.Now().UnixMilli()
	keys := make([]string, 0, len(s.data))
	for k, rec := range s.data {
		if rec.exp != -1 && rec.exp <= nowMs {
			delete(s.data, k)
			continue
		}
		keys = append(keys, k)
	}
	return keys
}

func (s *InMemoryStore) GetAllValues() [][]byte {
	nowMs := time.Now().UnixMilli()
	values := make([][]byte, 0, len(s.data))
	for k, rec := range s.data {
		if rec.exp != -1 && rec.exp <= nowMs {
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
