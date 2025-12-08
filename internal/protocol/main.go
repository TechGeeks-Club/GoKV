package protocol

import (
	"bufio"

	"github.com/B-AJ-Amar/gokv/internal/store"
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

const (
	ExpireNone = iota
	ExpireEX
	ExpirePX
	ExpireEXAT
	ExpirePXAT
	ExpireKEEPTTL
)

// todo : add set flags support
type RESPSetArgs struct {
	expType int8
	expVal  int
	xx      bool
	nx      bool
	keepTTL bool
}

type RESPReq struct {
	cmd     string
	argsLen int
	args    []string
	setArgs RESPSetArgs
}

type RESPRes struct {
	msgType int
	message string
}

type Protocol interface {
	Parse(msg string) (*RESPReq, error)
	Process(req *RESPReq, mem *store.InMemoryStore) (*RESPRes, error)
	Send(writer *bufio.Writer, res *RESPRes) (int, error)
	SendError(msg string)
}

type RESP struct{}
