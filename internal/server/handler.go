package server

import (
	"bufio"
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
		resp.SendError(w, err.Error())
		return
	}

	res, err := resp.Process(req, mem)
	if err != nil {
		resp.SendError(w, err.Error())
		return
	}

	resp.Send(w, res)

}
