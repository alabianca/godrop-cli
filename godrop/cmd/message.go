package cmd

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

type Header struct {
	mLength uint16
	sLength uint16
}

func (m *Header) Marshall() ([]byte, error) {
	header := make([]byte, 0)

	mlength, err := uint16ToBytes(m.mLength)

	if err != nil {
		return nil, err
	}

	slength, err := uint16ToBytes(m.sLength)

	if err != nil {
		return nil, err
	}

	header = append(header, mlength...)
	header = append(header, slength...)

	return header, nil

}

func newMessage(msg []byte, publicKey *rsa.PublicKey) (*Message, error) {
	label := []byte("")
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		publicKey,
		msg,
		label,
	)

	if err != nil {
		return nil, err
	}

	message := &Message{
		cipher: ciphertext,
		msg:    msg,
	}

	return message, nil
}

type Message struct {
	header    Header
	cipher    []byte
	signature []byte
	msg       []byte
}

func (m *Message) Sign(privateKey *rsa.PrivateKey) error {
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := m.msg
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(
		rand.Reader,
		privateKey,
		newhash,
		hashed,
		&opts,
	)

	if err != nil {
		return err
	}

	m.signature = signature
	m.header = Header{
		mLength: uint16(len(m.cipher)),
		sLength: uint16(len(m.signature)),
	}

	return nil
}

func (m *Message) Marshall() ([]byte, error) {
	header, err := m.header.Marshall()

	if err != nil {
		return nil, err
	}

	message := make([]byte, 0)
	message = append(message, MESSAGE)
	message = append(message, header...)
	message = append(message, m.cipher...)
	message = append(message, m.signature...)
	message = append(message, END_OF_TEXT)

	return message, nil
}

func (m *Message) UnMarshall(msg []byte) error {
	fmt.Println("Unmarshalling")
	fmt.Println(msg)
	header := Header{
		mLength: binary.BigEndian.Uint16(msg[1:3]),
		sLength: binary.BigEndian.Uint16(msg[3:5]),
	}

	fmt.Printf("%d %d\n", header.mLength, header.sLength)

	startOfCipher := 5
	endOfCipher := 5 + header.mLength
	startOfSignature := endOfCipher
	endOfSignature := startOfSignature + header.sLength

	m.header = header

	m.cipher = msg[startOfCipher:endOfCipher]
	m.signature = msg[startOfSignature:endOfSignature]

	return nil
}

func (m *Message) Decrypt(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) error {
	hash := sha256.New()
	label := []byte("")

	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, m.cipher, label)

	if err != nil {
		return err
	}

	// test signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := plainText
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	err = rsa.VerifyPSS(
		publicKey,
		newhash,
		hashed,
		m.signature,
		&opts,
	)

	if err != nil {
		return err
	}

	fmt.Println(string(plainText))
	m.msg = plainText

	return nil
}
