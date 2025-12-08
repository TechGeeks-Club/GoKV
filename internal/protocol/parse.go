package protocol

import (
	"bufio"
	"strconv"
	"strings"
)

var (
	allowedCommands = [...]string{"set", "get", "del", "incr", "exists", "ping", "hello"}
)

func (r *RESP) Parse(reader *bufio.Reader) (*RESPReq, error) {
	req := RESPReq{}
	msg, err := reader.ReadString('\n')
	argLen := 0
	if err != nil {
		return nil, ErrInvalidFormat
	}
	msg = strings.TrimRight(msg, "\r\n")

	if msg[0] == '*' {
		req.argsLen, err = strconv.Atoi(msg[1:])
		if err != nil {
			return nil, ErrParseLen
		}
	} else {
		return nil, ErrInvalidFormat
	}

	for i := 0; i < req.argsLen; i++ {
		// "$N" len of the next arg
		msg, err := reader.ReadString('\n')
		if err != nil {
			return nil, ErrInvalidFormat
		}
		msg = strings.TrimRight(msg, "\r\n")
		if msg[0] == '$' {
			argLen, err = strconv.Atoi(msg[1:])
			if err != nil {
				return nil, ErrParseLen
			}
		} else {
			return nil, ErrInvalidFormat
		}
		// the arg
		msgArg, err := reader.ReadString('\n')
		if err != nil {
			return nil, ErrInvalidFormat
		}
		msgArg = strings.TrimRight(msgArg, "\r\n")

		if argLen != len(msgArg) {
			return nil, ErrWrongArgLen
		}
		req.args = append(req.args, msgArg)

	}
	cmd := strings.ToLower(req.args[0])

	switch cmd {
	case "set":

		if len(req.args) < 3 || len(req.args) > 7 {
			return nil, ErrWrongNumberArgs
		}
		// check from the expiry and XX, NX args
		if len(req.args) > 3 {
			i := 3
			for i < len(req.args) {
				arg := strings.ToUpper(req.args[i])
				if arg == "EX" || arg == "PX" || arg == "EXAT" || arg == "PXAT" {
					if i+1 >= len(req.args) {
						return nil, ErrWrongNumberArgs
					}
					_, err := strconv.Atoi(req.args[i+1])
					if err != nil {
						return nil, ErrInvalidExpireTime
					}
					i += 2
				} else if arg == "NX" || arg == "XX" || arg == "KEEPTTL" {
					i++
				} else {
					return nil, ErrWrongNumberArgs
				}
			}
		}
	case "get":
		if len(req.args) != 2 {
			return nil, ErrWrongNumberArgs
		}
	case "del":
		if len(req.args) < 2 {
			return nil, ErrWrongNumberArgs
		}
	case "exists":
		if len(req.args) < 2 {
			return nil, ErrWrongNumberArgs
		}
	case "ping":
		if len(req.args) != 1 {
			return nil, ErrWrongNumberArgs
		}
	default:
		return nil, ErrUnknownCommand
	}

	req.cmd = cmd

	return &req, nil
}
