
package server

import(
	"bufio"
)

func HandleConnection(conn net.Conn){
	defer conn.Close()

	r := bufio.NewReader(conn)
    w := bufio.NewWriter(conn)

	for {
        if err != nil {
            fmt.Println("Error reading:", err)
            return
        }

		// TODO : parse , (if there is an error in parsing that can mean that this is a part of the prev request (if the prev request ends with full buff))
		// TODO : receive
	}
}