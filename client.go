package main

import (
	"bufio"
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
	//Open the file
	file, err := os.Open(strings.TrimSpace(fileName))
	if err != nil {
		fmt.Println(err)
		return
	}

	// make sure to close the file
	defer file.Close()

	// Convert string name to bytes

	fileNameInBytes := []byte(fileName)
	fileNameInBytesSize := len(fileNameInBytes)

	// Send file name

	// Send the file to the connection
	n, err := io.Copy(connection, file)
	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	fmt.Println(n, "bytes sent")
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
