
package server

func HandleConnection(conn net.Conn){
	defer conn.Close()
	bufferLen := 4096
	isFullBuffer := false
    buffer := make([]byte, bufferLen)
	for {
		n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Error reading:", err)
            return
        }
		if n == bufferLen {
			isFullBuffer = true
		}
		// TODO : parse , (if there is an error in parsing that can mean that this is a part of the prev request (if the prev request ends with full buff))
		// TODO : receive
	}
}