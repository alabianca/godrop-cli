package main

import (
	"bufio"
	"io"
	"os"

	"github.com/alabianca/godrop-cli/godrop/cmd"

	"github.com/alabianca/godrop"
)

func main() {

	cmd.RootCmd.Execute()
	// drop, err := godrop.NewGodrop()

	// if err != nil {
	// 	panic(err)
	// }

	// connStrategy := drop.NewP2PConn("mdns")

	// p2pConn, err := connStrategy.Connect("")

	// if err != nil {
	// 	fmt.Println("Could Not establish P2P Connection")
	// 	os.Exit(1)
	// }

	// pipe(readStdin(), readFromPeer(p2pConn), p2pConn)

}

func readStdin() chan []byte {

	quitChan := make(chan []byte)

	go func(quit chan []byte) {
		buf := make([]byte, 1024)
		reader := bufio.NewReader(os.Stdin)
		result := make([]byte, 0)

		for {
			n, err := reader.Read(buf)

			if err != nil {
				if err == io.EOF {
					result = append(result, buf[:n]...)
					quit <- result
					close(quit)
					return
				}
			}
			result = append(result, buf[:n]...)
		}
	}(quitChan)

	return quitChan
}

func readFromPeer(conn *godrop.P2PConn) chan []byte {
	result := make(chan []byte)
	go func(quit chan []byte) {
		buf := make([]byte, 1024)
		data := make([]byte, 0)

		for {
			n, err := conn.Read(buf)

			if err != nil {
				if err == io.EOF {
					data = append(data, buf[:n]...)
					result <- data

				}
				close(result)
				return
			}

			data = append(data, buf[:n]...)
		}

	}(result)

	return result
}

func pipe(stdinChan, fromPeerChan <-chan []byte, conn *godrop.P2PConn) {

	for {
		select {
		case data := <-stdinChan:
			//write to the peer
			writeToPeer(data, conn)
			return

		case data := <-fromPeerChan:
			//write to stdout
			writeToStdout(data)
			return
		}
	}
}

func writeToPeer(data []byte, conn *godrop.P2PConn) {
	conn.Write(data)
	conn.Close()
}

func writeToStdout(data []byte) {
	writer := bufio.NewWriter(os.Stdout)
	writer.Write(data)
	writer.Flush()
}
