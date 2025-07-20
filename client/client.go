package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	mode := os.Args[1]
	switch mode {
	case "-r":
		if len(os.Args) != 3 {
			printUsage()
			os.Exit(1)
		}
		receiveMessages(os.Args[2])
	case "-s":
		if len(os.Args) != 5 {
			printUsage()
			os.Exit(1)
		}
		sendMessage(os.Args[2], os.Args[3], os.Args[4])
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: blitz -s <username> <target> \"<message>\"")
	fmt.Println("Or: blitz -r <username>")
}

func receiveMessages(username string) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT, config)
	if err != nil {
		log.Printf("Failed to connect to sever: %v\n", err)
		return
	}
	fmt.Fprintln(conn, username)
	fmt.Fprintln(conn, "server")
	fmt.Fprintln(conn, "false")
	fmt.Fprintln(conn, "register")
	fmt.Fprintln(conn, "0")

	fmt.Printf("Registered as %s. Listening for messages...\n", username)

	reader := bufio.NewReader(conn)
	input := bufio.NewReader(os.Stdin)
	for {

		sender, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read sender: %v\n", err)
			return
		}
		reader.ReadString('\n') //Discard the target
		isFileStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read isFile flag: %v\n", err)
			continue
		}
		isFile, err := strconv.ParseBool(strings.TrimSpace(isFileStr))
		if err != nil {
			fmt.Printf("Failed to parse isFile flag: %v\n", err)
			return
		}
		fileName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read filename: %v\n", err)
			return
		}
		fileSizeStr, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read filesize: %v\n", err)
			return
		}
		fileSize, err := strconv.ParseInt(strings.TrimSpace(fileSizeStr), 10, 64)
		if err != nil {
			fmt.Printf("Failed to parse filesize: %v\n", err)
			return
		}

		sender = strings.TrimSpace(sender)
		fileName = strings.TrimSpace(fileName)
		fmt.Printf("File from %s: %s. Accept [y/n]\n", sender, fileName)
		if line, _ := input.ReadString('\n'); strings.TrimSpace(line) == "y" {
			if isFile {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					log.Printf("Could not find user's home directory: %v. Saving to current directory.", err)
					homeDir = "." // "." means current directory
				}

				downloadPath := filepath.Join(homeDir, "Downloads", fileName)

				fmt.Printf("Saving file to: %s\n", downloadPath)
				file, err := os.Create(downloadPath)
				if err != nil {
					fmt.Printf("Failed to create file: %v\n", err)
					io.CopyN(io.Discard, reader, int64(fileSize))
					return
				}
				defer file.Close()
				bytesWritten, err := io.CopyN(file, reader, int64(fileSize))
				if err != nil {
					log.Printf("Error during file download: %v. %d bytes written\n", err, bytesWritten)
					return
				}
				fmt.Printf("File downloaded successfully. %d bytes were written.\n", bytesWritten)
			}
		} else {
			fmt.Println("File transfer declined.")
			io.CopyN(io.Discard, reader, int64(fileSize))
		}
	}
}

func sendMessage(username string, target string, filePath string) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT, config)
	if err != nil {
		log.Fatalf("Failed to connect to sever: %v", err)
	}
	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %v\n", err)
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Failed to get file info: %v\n", err)
		return
	}

	fileSize := fileInfo.Size()
	fileName := fileInfo.Name()

	writer := bufio.NewWriter(conn)

	writer.WriteString(username + "\n")
	writer.WriteString(target + "\n")
	writer.WriteString("true\n")
	writer.WriteString(fileName + "\n")
	writer.WriteString(strconv.FormatInt(fileSize, 10) + "\n")

	writer.Flush()

	bytesWritten, err := io.Copy(writer, file)
	if err != nil {
		log.Printf("Failed to send file: %v", err)
		return
	}

	fmt.Printf("File sent successfully. %d bytes sent to %s\n", bytesWritten, target)

	writer.Flush()

	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Server response: " + response)
}
