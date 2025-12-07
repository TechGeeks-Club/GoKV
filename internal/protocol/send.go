package protocol

import (
	"bufio"
	"errors"
	"fmt"
)

func (r *RESP) Send(writer *bufio.Writer, res *RESPRes) error {
	switch res.msgType {
	case SimpleRes:
		fmt.Fprintf(writer, "+%s\r\n", res.message)
	case ErrorRes:
		fmt.Fprintf(writer, "-%s\r\n", res.message)
	case BulkStrRes:
		fmt.Fprintf(writer, "*2\r\n$%d\r\n%s\r\n", len(res.message), res.message)
	case NotExistsRes:
		fmt.Fprintf(writer, "-1\r\n")
	case IntRes:
		fmt.Fprintf(writer, ":%s\r\n", res.message)
	case SpecialRes:
		fmt.Fprintf(writer, "%s", res.message)
	default:
		return errors.New("ERR unsupported msgType")
	}
	writer.Flush()
	return nil

}
func (r *RESP) SendError(writer *bufio.Writer, msg string) error {

	fmt.Fprintf(writer, "-%s\r\n", msg)
	writer.Flush()

	return nil

}
