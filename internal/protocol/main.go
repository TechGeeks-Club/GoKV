package protocol

import (
	"bufio"

	"github.com/B-AJ-Amar/gokv/internal/store"
)

var (
	allowedCommands = [...]string{"set", "get", "del", "incr", "incrby", "exists", "ping", "select", "hello"}
)

const (
	SimpleStrType   = '+'
	SimpleErrorType = '-'
	IntType         = ':'
	BulkStrType     = '$'
	ArrayType       = '*'
)

const (
	SimpleRes    = iota // +MESSAGE\r\n
	ErrorRes            // -ERROR\r\n
	BulkStrRes          // $n\r\nXXX\r\n
	NotExistsRes        // $-1\r\n
	IntRes              // :1\r\n
	SpecialRes          // to send directly hardcoded response
)

type RESPReq struct {
	cmd     string
	argsLen int
	args    []string
	setArgs store.SetArgs
}

type RESPRes struct {
	msgType int
	message string
}

type Protocol interface {
	Parse(reader *bufio.Reader, dbIndex *int) (*RESPReq, error)
	Process(req *RESPReq, mem *store.InMemoryStore) (*RESPRes, error)
	Send(writer *bufio.Writer, res *RESPRes) (int, error)
	SendError(msg string)
}

type RESP struct{}
