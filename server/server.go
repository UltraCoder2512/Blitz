package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Message struct {
	Sender  string
	Target  string
	Content string
	Conn    net.Conn
}

type FileTransfer struct {
	Sender   string
	Target   string
	IsFile   bool
	FileName string
	FileSize uint64
	Data     []byte
	Conn     net.Conn
}

type ListenerClient struct {
	Username string
	Conn     net.Conn
}

func main() {

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to load key pair: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", ":8080", config)

	if err != nil {
		log.Fatalf("Failed to open tcp server: %v", err)
	}

	defer listener.Close()

	fmt.Println("Blitz tcp server running on port: 8080.")

	msgChannel := make(chan FileTransfer)

	go dispatcher(msgChannel)

	for {

		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Failed to accept connection. Error: %v \n", err)
			continue
		}

		go handleConnection(conn, msgChannel)

	}
}

func handleConnection(conn net.Conn, msgChannel chan<- FileTransfer) {

	fmt.Printf("New client connected: %s \n", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)

	var err error

	defer func() {
		if err != nil {
			fmt.Printf("An error %v occured during execution. Closing connection and Exiting...", err)
			conn.Close()
		}
	}()

	sender, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to parse sender address. %v", err)
		return
	}
	target, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to parse target address. %v", err)
		return
	}
	isFileStr, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to parse isFile flag. %v", err)
		return
	}
	isFile, err := strconv.ParseBool(strings.TrimSpace(isFileStr))
	if err != nil {
		log.Printf("Failed to parse isFile flag to bool. %v", err)
		return
	}
	fileName, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to parse filename. %v", err)
		return
	}
	fileSizeStr, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to parse filesize. %v", err)
		return
	}
	fileSize, err := strconv.ParseUint(strings.TrimSpace(fileSizeStr), 10, 64)
	if err != nil {
		log.Printf("Failed to parse filesize to uint64. %v", err)
		return
	}
	data := make([]byte, fileSize)
	bytesRead, err := io.ReadFull(reader, data)
	if err != nil {
		log.Printf("Error reading file payload from %s (read %d bytes): %v", conn.RemoteAddr(), bytesRead, err)
		return
	}
	fileTransfer := FileTransfer{
		Sender:   strings.TrimSpace(sender),
		Target:   strings.TrimSpace(target),
		IsFile:   isFile,
		FileName: strings.TrimSpace(fileName),
		FileSize: fileSize,
		Data:     data,
		Conn:     conn,
	}
	msgChannel <- fileTransfer
}

func dispatcher(msgChan <-chan FileTransfer) {
	listeners := make(map[string]net.Conn)

	for {
		msg := <-msgChan
		if msg.Target == "server" {
			fmt.Printf("Listener Client %s registered. \n", msg.Sender)
			listeners[msg.Sender] = msg.Conn
		} else {
			fmt.Printf("Routing file from %s to %s \n", msg.Sender, msg.Target)
			targetConn, found := listeners[msg.Target]
			if found {
				writer := bufio.NewWriter(targetConn)
				writer.WriteString(msg.Sender + "\n")
				writer.WriteString(msg.Target + "\n")
				writer.WriteString(strconv.FormatBool(msg.IsFile) + "\n")
				writer.WriteString(msg.FileName + "\n")
				writer.WriteString(strconv.FormatUint(msg.FileSize, 10) + "\n")
				writer.Write(msg.Data)
				writer.Flush()
				fmt.Println("File transerred successfully!")
				fmt.Fprintln(msg.Conn, "File forwarded successfully!")
			} else {
				fmt.Printf("Target user %s not found.\n", msg.Target)
				fmt.Fprintln(msg.Conn, "Failed to forward file.")
			}
			msg.Conn.Close()
		}
	}
}
