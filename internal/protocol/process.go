package protocol

import (
	"strconv"

	"github.com/B-AJ-Amar/gokv/internal/store"
)

func (r *RESP) Process(req *RESPReq, mem *store.InMemoryStore) (*RESPRes, error) {
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
		mem.Set(req.args[1], []byte(req.args[2]))
		// todo : add exp
		response.msgType = SimpleRes
		response.message = "OK"
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
	case "ping":
		response.msgType = SimpleRes
		response.message = "PONG"
	default:
		response.msgType = ErrorRes
		response.message = "ERR unknown command"
	}

	return &response, nil
}
