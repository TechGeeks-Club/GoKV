package common

import "errors"

var (
	ErrInvalidFormat      = errors.New("ERR invalid format")
	ErrUnknownCommand     = errors.New("ERR unknown command")
	ErrWrongNumberArgs    = errors.New("ERR wrong number of arguments")
	ErrWrongArgLen        = errors.New("ERR wrong argument length")
	ErrKeyNotFound        = errors.New("ERR key not found")
	ErrParseLen           = errors.New("ERR parse len")
	ErrNotIntOROutOfRange = errors.New("ERR value is not an integer or out of range")
	ErrInvalidExpireTime  = errors.New("ERR invalid expire time")
	ErrInvalidIncrement   = errors.New("ERR invalid increment value")
	ErrInvalidDecrement   = errors.New("ERR invalid decrement value")
	ErrSyntaxError        = errors.New("ERR syntax error")
	ErrDBIndexOutOfRange  = errors.New("ERR DB index is out of range")
)
