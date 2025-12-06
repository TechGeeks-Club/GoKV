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
			response.msgType = ErrorRes
			response.message = ErrKeyNotFound.Error()
		} else {
			response.msgType = BulkStrRes
			response.message = string(res)
		}
	case "set":
		mem.Set(req.args[1], []byte(req.args[2]))
		// todo : add exp
		response.msgType = SpecialRes
		response.message = SetResMsg
	case "ping":
		response.msgType = SpecialRes
		response.message = PongResMsg
	case "hello":
		arg, err := strconv.Atoi(req.args[1])
		if err != nil {
			response.msgType = ErrorRes
			response.message = ErrKeyNotFound.Error()
		} else if arg == 2 {
			response.msgType = SpecialRes
			response.message = Hello2ResMsg
		} else {
			response.msgType = SpecialRes
			response.message = HelloResMsg
		}
	default:
		response.msgType = ErrorRes
		response.message = "ERR unknown command"
	}

	return &response, nil

}
