package protocol

import (
	"strconv"

	"github.com/B-AJ-Amar/gokv/internal/common"
	"github.com/B-AJ-Amar/gokv/internal/store"
)

func (r *RESP) Process(req *RESPReq, dbIndex *int, mem *store.InMemoryStore) (*RESPRes, error) {
	response := RESPRes{}
	switch req.cmd {
	case "get":
		res, err := mem.Get(req.args[1])
		if err == nil && res == nil {
			response.msgType = NotExistsRes
		} else {
			response.msgType = BulkStrRes
			response.message = string(res)
		}
	case "set":
		if req.argsLen == 3 {
			mem.Set(req.args[1], []byte(req.args[2]))
			response.msgType = SimpleRes
			response.message = "OK"
		} else {
			counter, oldRet, err := mem.Setx(req.args[1], []byte(req.args[2]), req.setArgs)
			if err != nil {
				response.msgType = ErrorRes
				response.message = "ERR syntax error"
			} else if counter == 0 {
				response.msgType = NotExistsRes
			} else if oldRet != nil {
				response.msgType = BulkStrRes
				response.message = string(oldRet)
			} else {
				response.msgType = SimpleRes
				response.message = "OK"
			}
		}

	case "del":
		deleted := mem.Del(req.args[1:])
		response.msgType = IntRes
		response.message = strconv.Itoa(deleted)
	case "exists":
		exists := mem.Exists(req.args[1:])
		response.msgType = IntRes
		response.message = strconv.Itoa(exists)
	case "incr":
		newVal, err := mem.Incrby(req.args[1], 1)
		if err != nil {
			response.msgType = ErrorRes
			response.message = "ERR value is not an integer or out of range"
		} else {
			response.msgType = IntRes
			response.message = strconv.Itoa(newVal)
		}
	case "incrby":
		by, _ := strconv.Atoi(req.args[2])
		newVal, err := mem.Incrby(req.args[1], by)
		if err != nil {
			response.msgType = ErrorRes
			response.message = "ERR value is not an integer or out of range"
		} else {
			response.msgType = IntRes
			response.message = strconv.Itoa(newVal)
		}
	case "decr":
		newVal, err := mem.Decrby(req.args[1], 1)
		if err != nil {
			response.msgType = ErrorRes
			response.message = "ERR value is not an integer or out of range"
		} else {
			response.msgType = IntRes
			response.message = strconv.Itoa(newVal)
		}
	case "decrby":
		by, _ := strconv.Atoi(req.args[2])
		newVal, err := mem.Decrby(req.args[1], by)
		if err != nil {
			response.msgType = ErrorRes
			response.message = "ERR value is not an integer or out of range"
		} else {
			response.msgType = IntRes
			response.message = strconv.Itoa(newVal)
		}

	case "ttl":
		ttl, err := mem.TTL(req.args[1])
		if err != nil {
			response.msgType = ErrorRes
			response.message = "ERR key does not exist"
		} else {
			response.msgType = IntRes
			response.message = strconv.Itoa(ttl)
		}
	case "expire":
		seconds, err := strconv.Atoi(req.args[2])
		if err != nil || seconds < 0 {
			return nil, common.ErrInvalidExpireTime
		}

		expired, err := mem.Expire(req.args[1], seconds)
		if err != nil {
			response.msgType = ErrorRes
			response.message = "ERR key does not exist"
		} else {
			response.msgType = IntRes
			response.message = strconv.Itoa(expired)
		}

	case "persist":
		persisted, err := mem.Persist(req.args[1])
		if err != nil {
			response.msgType = ErrorRes
			response.message = "ERR key does not exist"
		} else {
			response.msgType = IntRes
			response.message = strconv.Itoa(persisted)
		}

	case "select":
		newDBIndex, err := strconv.Atoi(req.args[1])
		if err != nil || newDBIndex < 0 || newDBIndex > common.MaxDBIndex {
			return nil, common.ErrNotIntOROutOfRange
		}

		*dbIndex = newDBIndex

		response.msgType = SimpleRes
		response.message = "OK"
	case "ping":
		response.msgType = SimpleRes
		response.message = "PONG"
	default:
		response.msgType = ErrorRes
		response.message = "ERR unknown command"
	}

	return &response, nil
}
