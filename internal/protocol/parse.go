package protocol

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/B-AJ-Amar/gokv/internal/common"
	"github.com/B-AJ-Amar/gokv/internal/store"
)

var (
	allowedCommands = [...]string{"set", "get", "del", "incr", "incrby", "exists", "ping", "hello"}
)

func getExpireType(arg string) int8 {
	switch arg {
	case "EX":
		return store.ExpireEX
	case "PX":
		return store.ExpirePX
	case "EXAT":
		return store.ExpireEXAT
	case "PXAT":
		return store.ExpirePXAT
	case "KEEPTTL":
		return store.ExpireKEEPTTL
	default:
		return store.ExpireNone
	}
}

func (r *RESP) Parse(reader *bufio.Reader) (*RESPReq, error) {
	req := RESPReq{}
	msg, err := reader.ReadString('\n')
	argLen := 0
	if err != nil {
		return nil, common.ErrInvalidFormat
	}
	msg = strings.TrimRight(msg, "\r\n")

	if msg[0] == '*' {
		req.argsLen, err = strconv.Atoi(msg[1:])
		if err != nil {
			return nil, common.ErrParseLen
		}
	} else {
		return nil, common.ErrInvalidFormat
	}

	for i := 0; i < req.argsLen; i++ {
		// "$N" len of the next arg
		msg, err := reader.ReadString('\n')
		if err != nil {
			return nil, common.ErrInvalidFormat
		}
		msg = strings.TrimRight(msg, "\r\n")
		if msg[0] == '$' {
			argLen, err = strconv.Atoi(msg[1:])
			if err != nil {
				return nil, common.ErrParseLen
			}
		} else {
			return nil, common.ErrInvalidFormat
		}
		// the arg
		msgArg, err := reader.ReadString('\n')
		if err != nil {
			return nil, common.ErrInvalidFormat
		}
		msgArg = strings.TrimRight(msgArg, "\r\n")

		if argLen != len(msgArg) {
			return nil, common.ErrWrongArgLen
		}
		req.args = append(req.args, msgArg)

	}
	cmd := strings.ToLower(req.args[0])

	switch cmd {
	case "set":

		if len(req.args) < 3 || len(req.args) > 7 {
			return nil, common.ErrWrongNumberArgs
		}
		// check from the expiry and XX, NX args
		if len(req.args) > 3 {
			i := 3
			for i < len(req.args) {
				arg := strings.ToUpper(req.args[i])
				if arg == "EX" || arg == "PX" || arg == "EXAT" || arg == "PXAT" {
					req.setArgs.ExpType = getExpireType(arg)
					if i+1 >= len(req.args) {
						return nil, common.ErrWrongNumberArgs
					}
					val, err := strconv.Atoi(req.args[i+1])
					if err != nil {
						return nil, common.ErrInvalidExpireTime
					}
					req.setArgs.ExpVal = val
					i += 2
				} else if arg == "NX" {
					req.setArgs.NX_XX = 1
					i++
				} else if arg == "XX" {
					req.setArgs.NX_XX = 2
					i++
				} else if arg == "KEEPTTL" {
					req.setArgs.KeepTTL = true
					i++
				} else if arg == "GET" {
					req.setArgs.Get = true
					i++
				} else {
					return nil, common.ErrWrongNumberArgs
				}
			}
		}
	case "get":
		if len(req.args) != 2 {
			return nil, common.ErrWrongNumberArgs
		}
	case "del":
		if len(req.args) < 2 {
			return nil, common.ErrWrongNumberArgs
		}
	case "exists":
		if len(req.args) < 2 {
			return nil, common.ErrWrongNumberArgs
		}
	case "incr":
		if len(req.args) != 2 {
			return nil, common.ErrWrongNumberArgs
		}
	case "incrby":
		if len(req.args) != 3 {
			return nil, common.ErrWrongNumberArgs
		}
		_, err := strconv.Atoi(req.args[2])
		if err != nil {
			return nil, common.ErrInvalidIncrement
		}
	case "ping":
		if len(req.args) != 1 {
			return nil, common.ErrWrongNumberArgs
		}
	default:
		return nil, common.ErrUnknownCommand
	}

	req.cmd = cmd

	return &req, nil
}
