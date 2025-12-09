package server

import (
	"bufio"
	"net"

	"github.com/B-AJ-Amar/gokv/internal/protocol"
	"github.com/B-AJ-Amar/gokv/internal/store"
)

func HandleConnection(conn net.Conn, mem *[]*store.InMemoryStore) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	dbIndex := 0

	for {
		resp := protocol.RESP{}
		req, err := resp.Parse(r)
		if err != nil {
			resp.SendError(w, err.Error())
			return
		}

		res, err := resp.Process(req, &dbIndex, (*mem)[dbIndex])
		if err != nil {
			resp.SendError(w, err.Error())
			return
		}

		resp.Send(w, res)

	}

}
