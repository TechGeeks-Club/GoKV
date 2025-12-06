package protocol

import (
	"fmt"

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
	SetResMsg    = "+OK\r\n"
	PongResMsg   = "+PONG\r\n"
	HelloResMsg  = "-NOPROTO unsupported protocol version\r\n"
	Hello2ResMsg = "*14\r\n$6\r\nserver\r\n$5\r\nGpKv\r\n$7\r\nversion\r\n$11\r\n0.0.0-alpha\r\n$5\r\nproto\r\n:2\r\n$2\r\nid\r\n:4\r\n$4\r\nmode\r\n$9\r\nstandalone\r\n$4\r\nrole\r\n$6\r\nmaster\r\n$7\r\nmodules\r\n*0\r\n"
)

const (
	SimpleRes    = iota // +MESSAGE\r\n
	ErrorRes            // -ERROR\r\n
	BulkStrRes          // $n\r\nXXX\r\n
	NotExistsRes        // $-1\r\n
	IntRes              // :1\r\n
	SpecialRes
)

type RESPReq struct {
	cmd     string
	argsLen int
	args    []string
}

type RESPRes struct {
	msgType int
	message string
}

type Protocol interface {
	Parse(msg string) (*RESPReq, error)
	Process(req *RESPReq, mem *store.InMemoryStore) (*RESPRes, error)
	Send()
	SendError(msg string)
}

type RESP struct{}

func main() {
	fmt.Println("Hello protocol")
}
