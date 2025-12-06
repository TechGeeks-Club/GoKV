package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/B-AJ-Amar/gokv/internal/protocol"
	"github.com/B-AJ-Amar/gokv/internal/store"
)

func HandleConnection(conn net.Conn, mem *store.InMemoryStore) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	resp := protocol.RESP{}

	req, err := resp.Parse(r)
	if err != nil {
		w.WriteString(fmt.Sprintf("-%s\r\n", err.Error()))
	}

	res, err := resp.Process(req, mem)
	// TODO : process
	// TODO : response

}
