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

var currentChannels = make([]string, 0)

func handleServerResponse(connection net.Conn) {
	// Buffer that holds command protocol number
	commandProtocolBuffer := make([]byte, 1)

	for {
		// Read a single byte, this byte determines the command (response) that
		// is made to this particular client
		// If that byte is not identified as a proper command response, all the
		// connection will be flushed leaving it empty. This will take care of
		// errors and unexpected behaviour

		// Step 1: Read command buffer
		_, err := io.ReadFull(io.LimitReader(connection, 1), commandProtocolBuffer)

		if err != nil {
			fmt.Println("Step 1 error:", err.Error())
			break
		}

		// Convert Command Buffer
		commandNumber := int(commandProtocolBuffer[0])
		fmt.Printf("The command number is %d\n", commandNumber)

		switch commandNumber {
		// Send File
		case 24:
			receiveFile(connection)
		case 34:
			//client.handleChannelSubscription()
		case 44:
			//client.changeUserName()
		default:
			fmt.Println("Unknown command received: Flushing connection")
			w := bufio.NewWriter(connection)
			w.Flush()
		}
	}
}

func handleUserCommand(userTextInput string, connection net.Conn) {
	trimmedText := strings.Split(strings.TrimSpace(userTextInput), " ")
	command := trimmedText[0]
	arguments := trimmedText[1:]

	switch command {
	case "SEND": // Command 24
		if len(arguments) == 0 {
			fmt.Println("Missing command arguments!")
		} else {
			sendFile(connection, arguments[0])
		}

	case "SUB": // Command 34
		if len(arguments) == 0 {
			fmt.Println("Missing command arguments!")
		} else {
			subscribeToChannel(connection, arguments[0])
		}

	case "USERNAME": // Command 44
		if len(arguments) == 0 {
			fmt.Println("Missing command arguments!")
		} else {
			changeUserName(connection, arguments[0])
		}

	case "CHNLS":
		listChannels()
	case "MSG":
		fmt.Println("TODO")
	case "HELP":
		fmt.Println("Available commands are: HELP SEND SUB USERNAME CHNLS MSG")
	default:
		fmt.Printf("Wrong Syntax, use the 'HELP' command if you need more info")
	}
}

func receiveFile(connection net.Conn) {
	fmt.Println("####Reading####")
	// The protocol number for receiving a file is 24

	// Buffer that holds file name length
	fileNameLengthBuffer := make([]byte, 4)

	// Buffer that holds file length
	fileLengthBuffer := make([]byte, 8)

	// This for should be removed*
	//for {
	// Step 1: Read name file length
	fmt.Println("Step 1: Read name file length")

	// Read File Name length buffer
	// bytesRead, err = io.ReadFull(client.connection, fileNameLengthBuffer)
	bytesRead, err := io.ReadFull(io.LimitReader(connection, 4), fileNameLengthBuffer)
	if err != nil {
		fmt.Println("Step 1 error:", err.Error())
		//break
	}
	fmt.Println(bytesRead)

	// Convert File Name length buffer to the length of the file name

	fmt.Printf("File name length buffer %v\n", fileNameLengthBuffer)

	fileNameLength := int32(binary.LittleEndian.Uint32(fileNameLengthBuffer))

	fmt.Printf("Step 2: Received name file length: %d\n", fileNameLength)

	// Step 3: Read file name
	fmt.Println("Step 3: Read file name")

	// Buffer that holds the name of the file
	fileNameBuffer := make([]byte, int64(fileNameLength))

	// Receive fileNameLength bytes which will be the name of the file
	fmt.Printf("We should receive exactly a string of bytes: %d \n", int64(fileNameLength))

	bytesRead, err = io.ReadFull(io.LimitReader(connection, int64(fileNameLength)), fileNameBuffer)
	if err != nil {
		fmt.Println("Step 3 error:", err.Error())
		//break
	}

	// Convert FileName buffer to string
	fileName := string(fileNameBuffer)
	fmt.Printf("Step 3: The name of the file is: %s\n", fileName)

	// Step 4: Get the buffer size of the file
	bytesRead, err = io.ReadFull(io.LimitReader(connection, 8), fileLengthBuffer)
	if err != nil {
		fmt.Println("Step 4: Error reading:", err.Error())
		//break
	}

	// Convert the buffer size slice to a number
	fileLength := int64(binary.LittleEndian.Uint64(fileLengthBuffer))

	fmt.Printf("Step 4: We should receive a file of size: %d bytes\n", fileLength)

	// Step 5
	// Read the buffer and copy it to the created file with name fileName
	// Create the file

	receivedFile, err := os.Create(fileName)
	defer receivedFile.Close()

	fmt.Println("Copying...")

	// Read the file and copy it into fileName
	bytes, err := io.CopyN(receivedFile, connection, fileLength)
	if err != nil {
		fmt.Println("Step 5: Error reading:", err.Error())
		//break
	}
	fmt.Printf("Step 5: Bytes read: %d", bytes)
	fmt.Println("Copied successfully")

	if err != nil {
		fmt.Println("Step 5: Error reading:", err.Error())
		//break
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

// Subscribe to the channel named channelName, if the channel does not exist,
// it will be created and then subscribed to it
func subscribeToChannel(connection net.Conn, channelName string) {

	// Step 1: Send command
	// The protocol number for sending a file is 34
	_, err := connection.Write([]byte{34})
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("Command byte sent: Sent %d bytes", n)

	// Step 2: Send channel Name length
	// Convert string name to bytes and get the length
	channelNameInBytes := []byte(channelName)
	channelNameBytesSize := len(channelNameInBytes)
	//fmt.Printf("Step 2: the size in bytes of the channel names length is %d bytes\n", channelNameBytesSize)

	channelNameBufferLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(channelNameBufferLength, uint32(channelNameBytesSize))
	//fmt.Printf("The buffer of the length is %v \n", channelNameBufferLength)

	// Step 2: Send file name size
	_, err = connection.Write(channelNameBufferLength)

	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	//fmt.Printf("Step 2: sent channel name length %d bytes\n", n)

	//Step 3: Send channel name
	_, err = connection.Write(channelNameInBytes)

	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	//fmt.Printf("Step 3: sent channel name %d bytes\n", n)
	currentChannels = append(currentChannels, channelName)
	fmt.Printf("Channels subscribed to: %v \n", currentChannels)
}

// Changes current client username to 'newUserName'
func changeUserName(connection net.Conn, newUsername string) {
	// Step 1: Send command
	// The protocol number for changing the username is 44
	_, err := connection.Write([]byte{44})
	if err != nil {
		fmt.Println("We couldn't change the Username")
		return
	}

	// Send 'newUserName' length
	// Convert string name to bytes and get the length
	newUsernameBuffer := []byte(newUsername)
	newUsernameLength := len(newUsernameBuffer)

	newUsernameBufferLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(newUsernameBufferLength, uint32(newUsernameLength))

	// Send new username name size
	_, err = connection.Write(newUsernameBufferLength)

	if err != nil {
		fmt.Println("We couldn't change the Username")
		fmt.Println(err)
		return
	}

	// Send new username buffer
	bytesSent, err := connection.Write(newUsernameBuffer)
	if err != nil || bytesSent != newUsernameLength {
		fmt.Println("We couldn't change the Username")
		fmt.Println(err)
		return
	}

	fmt.Printf("Changed username to: %s\n", newUsername)
}

func listChannels() {
	fmt.Printf("Channels subscribed to: %v \n", currentChannels)
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
	go handleServerResponse(serverConnection)

	for {
		// Read user input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		userInputText, _ := reader.ReadString('\n')

		handleUserCommand(userInputText, serverConnection)
	}
}
