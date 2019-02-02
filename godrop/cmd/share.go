package cmd

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/alabianca/godrop"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share [INSTANCE]",
	Short: "Share a file via mdns or tcp holepunch",
	Long: `Share a file via mdns or tcp holepunch. 
	Run godrop share mdns to share a file via mdns.
	Run godrop share hp to share a file with a remote peer via tcp holepunch.`,
	Run: runShare,
}

func runShare(command *cobra.Command, args []string) {

	if len(args) <= 0 {
		command.Usage()
		os.Exit(1)
	}

	peer := args[0]

	share, err := newShare()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	share.Connect(peer)

}

func init() {
	RootCmd.AddCommand(shareCmd)

}

func newShare() (*Share, error) {
	prvKey, err := loadPrivateKey()

	if err != nil {
		return nil, err
	}

	pubKey, err := loadPublicKey()

	if err != nil {
		return nil, err
	}

	drop, err := configGodropMdns()

	if err != nil {
		return nil, err
	}

	share := &Share{
		drop:         drop,
		myPublicKey:  pubKey,
		myPrivateKey: prvKey,
		sharePath:    "",
		wg:           new(sync.WaitGroup),
	}

	return share, nil

}

type Share struct {
	drop          *godrop.Godrop
	myPublicKey   *rsa.PublicKey
	myPrivateKey  *rsa.PrivateKey
	peerPublicKey *rsa.PublicKey
	reader        *bufio.Reader
	writer        *bufio.Writer
	sharePath     string
	wg            *sync.WaitGroup
}

func (s *Share) Connect(instance string) error {
	conn, err := s.drop.Connect(instance)

	if err != nil {
		return err
	}

	s.reader = bufio.NewReader(conn)
	s.writer = bufio.NewWriter(conn)

	if err := s.handshake(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Waiting...")
	s.wg.Wait()

	return nil

}

// start the godrop handshake by sending your public key
func (s *Share) handshake() error {
	authPacket := make([]byte, 0)
	pubKey := x509.MarshalPKCS1PublicKey(s.myPublicKey)

	authPacket = append(authPacket, byte(HANDSHAKE))
	authPacket = append(authPacket, pubKey...)
	authPacket = append(authPacket, byte(END_OF_TEXT))

	_, err := s.writer.Write(authPacket)

	if err != nil {
		return err
	}

	if err := s.writer.Flush(); err != nil {
		return err
	}

	buf := make([]byte, 272)
	response := make([]byte, 0)
	totalRead := 0
	for {
		n, err := s.reader.Read(buf)

		if err != nil {
			if totalRead == 0 {

				return fmt.Errorf("Access Denied")
			}

			if err == io.EOF {
				//we have the peers public key
				return fmt.Errorf("Connection got closed")
			}
		} else {
			response = append(response, buf[:n]...)

			if len(response) == 272 {
				break
			}
		}

		//fmt.Println(buf[:n])
	}

	peerKeyBytes := response[1:271]
	peerKey, err := x509.ParsePKCS1PublicKey(peerKeyBytes)

	if err != nil {
		return err
	}

	s.peerPublicKey = peerKey

	return nil

}

// responde to a handshake message by responding with your public key
func (s *Share) handleHandshake(handshakeMsg []byte) error {
	peerKeyBytes := handshakeMsg[1:271]

	peerKey, err := x509.ParsePKCS1PublicKey(peerKeyBytes)

	if err != nil {
		return err
	}

	s.peerPublicKey = peerKey

	authPacket := make([]byte, 0)
	pubKey := x509.MarshalPKCS1PublicKey(s.myPublicKey)

	authPacket = append(authPacket, byte(HANDSHAKE_ACK))
	authPacket = append(authPacket, pubKey...)
	authPacket = append(authPacket, byte(END_OF_TEXT))

	n, err := s.writer.Write(authPacket)

	if err != nil || n < len(authPacket) {
		return err
	}

	if err := s.writer.Flush(); err != nil {
		return err
	}

	return nil

}
