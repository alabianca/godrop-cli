package cmd

import (
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/alabianca/godrop"

	"github.com/spf13/cobra"
)

const (
	BUF_SIZE = 1024
)

var advertiseCmd = &cobra.Command{
	Use:   "advertise [FILE_TO_SHARE]",
	Short: "Share a file with peers in the local network",
	Run:   runAdvertise,
}

func runAdvertise(command *cobra.Command, args []string) {

	if len(args) < 1 {
		command.Usage()
		os.Exit(1)
	}

	fPath := args[0]

	drop, err := configGodropMdns()

	if err != nil {
		log.Fatal(err)
	}

	advertiser, err := NewAdvertiser(fPath, drop)

	defer advertiser.Shutdown()

	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("Shutdown...")
	}

}

func init() {
	RootCmd.AddCommand(advertiseCmd)

}

type Advertiser struct {
	sharePath string
	server    *godrop.Server
	shutDown  chan int
	wg        *sync.WaitGroup
}

func (ad *Advertiser) Shutdown() {
	log.Println("shutting down advertiser...")
	ad.server.Shutdown()
	ad.shutDown <- 1
	ad.wg.Wait()
}

func NewAdvertiser(sharePath string, drop *godrop.Godrop) (*Advertiser, error) {
	if err := checkFile(sharePath); err != nil {
		return nil, err
	}

	server, err := drop.NewMDNSService(sharePath)

	if err != nil {
		return nil, err
	}

	ad := &Advertiser{
		sharePath: sharePath,
		server:    server,
		shutDown:  make(chan int),
		wg:        new(sync.WaitGroup),
	}

	go acceptLoop(ad.sharePath, ad.wg, ad.server, ad.shutDown)

	return ad, nil
}

func acceptLoop(dir string, wg *sync.WaitGroup, server *godrop.Server, shutdown chan int) {
	wg.Add(1)
	sessionChan := make(chan *godrop.Session)
	go func(s chan *godrop.Session) {
		for {
			sesh, err := server.Accept()

			if err != nil {
				log.Println("could not accept connection ", err)
				continue
			}

			s <- sesh
		}
	}(sessionChan)

	for {
		select {
		case <-shutdown:
			log.Println("Shutting down accept loop")
			wg.Done()
			return
		case sesh := <-sessionChan:
			go transferDir(dir, sesh)
		}
	}
}

func checkFile(fPath string) error {
	_, err := os.Stat(fPath)

	if os.IsNotExist(err) {
		return err
	}

	return nil
}

func transferDir(dir string, s *godrop.Session) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		writeHeader(s, path, info)

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		defer file.Close()

		if err != nil {
			return err
		}

		//write file
		buf := make([]byte, BUF_SIZE)

		for {
			n, err := file.Read(buf)

			if err != nil {
				if err == io.EOF {
					log.Println("Transferred: ", info.Name())
					break
				} else {
					log.Println("Error reading: ", info.Name())
				}
			}

			s.Write(buf[:n])
			s.Flush()
		}

		return nil

	})

	return err
}

func writeHeader(s *godrop.Session, path string, info os.FileInfo) {
	header := godrop.Header{
		Name:  info.Name(),
		Size:  info.Size(),
		IsDir: info.IsDir(),
		Path:  path,
	}

	s.WriteHeader(header)
}
