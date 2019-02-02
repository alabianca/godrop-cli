package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/alabianca/godrop"

	"github.com/spf13/cobra"
)

var advertiseCmd = &cobra.Command{
	Use:   "advertise [COMMAND]",
	Short: "Accept Files Shared by Peers",
	Run:   runAdvertise,
}

func runAdvertise(command *cobra.Command, args []string) {
	drop, err := configGodropMdns()

	if err != nil {
		log.Fatal(err)
	}

	server, err := drop.NewMDNSService()

	if err != nil {
		log.Fatal(err)
	}

	go mainLoop(server)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("Shutdown...")
		server.Shutdown()
	}

}

func init() {
	RootCmd.AddCommand(advertiseCmd)

}

func mainLoop(s *godrop.Server) {
	l, err := s.Listen()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if ok := acceptConnection(conn); !ok {
			fmt.Println("Closing connection")
			conn.Close()
			continue
		}

		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

	share, err := newShare()

	if err != nil {
		return
	}

	share.reader = bufio.NewReader(conn)
	share.writer = bufio.NewWriter(conn)

	go advertiseReadLoop(share)

	share.wg.Wait()
}

func advertiseReadLoop(s *Share) {
	s.wg.Add(1)

	buf := make([]byte, 1024)
	for {
		_, err := s.reader.Read(buf)

		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				s.wg.Done()
				break
			}
		}

		opType := int(buf[0])

		switch opType {
		case HANDSHAKE:
			if err := s.handleHandshake(buf); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func acceptConnection(conn net.Conn) bool {
	fmt.Printf("A peer from %s would like to connect. Allow connection? (y/n): ", conn.RemoteAddr().String())

	reader := bufio.NewReader(os.Stdin)

	for {
		answ, err := reader.ReadString('\n')

		if err != nil {
			continue
		}

		answ = strings.TrimSpace(strings.ToLower(answ))

		switch answ {
		case "y":
			return true
		case "n":
			return false
		default:
			continue
		}
	}
}
