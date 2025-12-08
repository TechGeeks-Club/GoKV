package protocol

import "errors"

var (
	ErrInvalidFormat      = errors.New("ERR invalid format")
	ErrUnknownCommand     = errors.New("ERR unknown command")
	ErrWrongNumberArgs    = errors.New("ERR wrong number of arguments")
	ErrWrongArgLen        = errors.New("ERR wrong argument lenght")
	ErrKeyNotFound        = errors.New("ERR key not found")
	ErrParseLen           = errors.New("ERR parse len")
	ErrNotIntOROutOfRange = errors.New("ERR value is not an integer or out of range")
	ErrInvalidExpireTime  = errors.New("ERR invalid expire time")
)
