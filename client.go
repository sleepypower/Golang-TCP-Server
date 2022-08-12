package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func handleUserCommand(userTextInput string, connection net.Conn) {
	trimmedText := strings.Split(strings.TrimSpace(userTextInput), " ")
	command := trimmedText[0]
	arguments := trimmedText[1:]

	println("#######")
	println(len(arguments))
	println("#######")

	if len(arguments) == 0 {
		fmt.Println("Missing command arguments!")
		return
	}
	switch command {
	case "SEND":
		sendFile(connection, arguments[0])
	case "SUB":
		fmt.Println("What channel would you like to subscribe to")
	case "CHNLS":
		fmt.Println("W")
	case "MSG":
		fmt.Println("W")
	case "HELP":
		fmt.Println("W")
	default:
		fmt.Printf("Wrong Syntax, use the 'HELP' command if you need more info")
	}
}

// The file named 'fileName' will be sent to all the clients subscribed with
// the same channels as the current client
func sendFile(connection net.Conn, fileName string) {

	// Check if filename exceeds 64 bytes

	// Step 1: Send command
	// The protocol number for sending a file is 24
	n, err := connection.Write([]byte{24})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Step 1: sent %d bytes\n", n)

	//Open the file
	file, err := os.Open(strings.TrimSpace(fileName))
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get file stats
	fileSt, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	// Get file size
	fileSize := fileSt.Size()

	// make sure to close the file
	defer file.Close()

	// Convert string name to bytes and get the length

	fileNameInBytes := []byte(fileName)
	fileNameInBytesSize := len(fileNameInBytes)
	fmt.Printf("Step 2: the size in bytes of the file names length is %d bytes\n", fileNameInBytesSize)

	fileNameBufferLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(fileNameBufferLength, uint32(fileNameInBytesSize))
	fmt.Printf("The buffer of the length is %v \n", fileNameBufferLength)

	// Step 2: Send file name size
	n, err = connection.Write(fileNameBufferLength)

	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	fmt.Printf("Step 2: sent %d bytes\n", n)

	// Step 3: Send file name buffer
	n, err = connection.Write([]byte(fileName))
	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	fmt.Printf("Step 3: sent %d bytes\n", n)

	// Step 4: Send file buffer size
	fmt.Println("file size in bytes is", fileSize)

	fileSizeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(fileSizeBuffer, uint64(fileSize))

	n, err = connection.Write([]byte(fileSizeBuffer))
	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	// Step 5: Send file buffer

	bytesWritten, err := io.Copy(connection, file)
	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}
	fmt.Printf("Sent %d bytes of the file named %s \n", bytesWritten, fileName)

}

func main() {

	// Connect with the server
	serverConnection, err := net.Dial("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	// // make sure to close the connection with the server
	defer serverConnection.Close()

	for {
		// Read user input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		userInputText, _ := reader.ReadString('\n')

		handleUserCommand(userInputText, serverConnection)

		// // Send user input message to the server
		// fmt.Fprintf(serverConnection, userInputText+"\n")

		// // Read server response
		// serverMessageResponse, _ := bufio.NewReader(serverConnection).ReadString('\n')

		// // Print server's response
		// fmt.Print("->: " + serverMessageResponse)

		// if strings.TrimSpace(string(userInputText)) == "STOP" {
		// 	fmt.Println("TCP client exiting...")
		// 	return
		// }
	}
}
