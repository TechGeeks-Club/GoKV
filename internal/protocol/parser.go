package protocol

import (
	"bufio"
	"strconv"
	"strings"
)

var (
	allowedCommands = [...]string{"set", "get", "ping", "hello"}
)

func (r *RESP) Parse(reader *bufio.Reader) (*RESPReq, error) {
	req := RESPReq{}
	msg, err := reader.ReadString('\n')
	argLen := 0
	if err != nil {
		return nil, ErrInvalidFormat
	}
	msg = msg[:len(msg)-1] // del \r

	if msg[0] == '*' {
		req.argsLen, err = strconv.Atoi(msg[1:])
		if err != nil {
			return nil, ErrInvalidFormat
		}
	} else {
		return nil, ErrInvalidFormat
	}

	for i := 0; i < req.argsLen; i++ {
		// "$N" len of the next arg
		msg, err := reader.ReadString('\n')
		msg = msg[:len(msg)-1] // del \r
		if err != nil {
			return nil, ErrInvalidFormat
		}
		if msg[0] == '$' {
			argLen, err = strconv.Atoi(msg[1:])
			if err != nil {
				return nil, ErrInvalidFormat
			}
		} else {
			return nil, ErrInvalidFormat
		}
		// the arg
		msgArg, err := reader.ReadString('\n')
		if err != nil {
			return nil, ErrInvalidFormat
		}
		msgArg = msgArg[:len(msg)-1] // del \r

		if argLen != len(msgArg) {
			return nil, ErrWrongArgLen

		}
		req.args = append(req.args, msgArg)

	}
	cmd := strings.ToLower(req.args[0])
	isValidCommand := false
	for _, v := range allowedCommands {
		if cmd == v {
			isValidCommand = true
			req.cmd = v
			break
		}

	}
	if !isValidCommand {
		return nil, ErrUnknownCommand
	}

	return &req, nil
}
