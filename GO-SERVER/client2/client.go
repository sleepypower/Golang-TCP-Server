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

// handles the server's response by identifying the first byte of the response
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
			subscribeToChannelResponse(connection)
		case 44:
			//client.changeUserName()
		default:
			fmt.Println("Unknown command received: Flushing connection")
			w := bufio.NewWriter(connection)
			w.Flush()
		}
	}
}

// Handles the user input command by triming and spliting the input and then
// defining what is the command that the user meant to use
func handleUserCommand(userTextInput string, connection net.Conn) {
	trimmedText := strings.Split(strings.TrimSpace(userTextInput), " ")
	command := trimmedText[0]
	arguments := trimmedText[1:]

	switch command {
	case "SEND": // Command 24
		if len(arguments) == 0 {
			fmt.Println("Missing command arguments!")
		} else if len(arguments) == 1 {
			println("Sending files to all subscribed channels!")
			sendFile(connection, arguments[0], "")
		} else {
			sendFile(connection, arguments[0], arguments[1])
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
		fmt.Println("Wrong Syntax, use the 'HELP' command if you need more info")
	}
}

// Receives a file sent by the server, including the size of the file, its name
// and length. Then saves the file in the directory
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

	//fmt.Printf("File name length buffer %v\n", fileNameLengthBuffer)

	fileNameLength := int32(binary.LittleEndian.Uint32(fileNameLengthBuffer))

	//fmt.Printf("Step 2: Received name file length: %d\n", fileNameLength)

	// Step 3: Read file name
	//fmt.Println("Step 3: Read file name")

	// Buffer that holds the name of the file
	fileNameBuffer := make([]byte, int64(fileNameLength))

	// Receive fileNameLength bytes which will be the name of the file
	//fmt.Printf("We should receive exactly a string of bytes: %d \n", int64(fileNameLength))

	bytesRead, err = io.ReadFull(io.LimitReader(connection, int64(fileNameLength)), fileNameBuffer)
	if err != nil {
		fmt.Println("Step 3 error:", err.Error())
		//break
	}

	// Convert FileName buffer to string
	fileName := string(fileNameBuffer)
	fmt.Printf("Receiving a file named: %s\n", fileName)

	// Step 4: Get the buffer size of the file
	bytesRead, err = io.ReadFull(io.LimitReader(connection, 8), fileLengthBuffer)
	if err != nil {
		fmt.Println("Step 4: Error reading:", err.Error())
		//break
	}

	// Convert the buffer size slice to a number
	fileLength := int64(binary.LittleEndian.Uint64(fileLengthBuffer))

	//fmt.Printf("Step 4: We should receive a file of size: %d bytes\n", fileLength)

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

	fmt.Printf("Bytes read: %d\n", bytes)

	if err != nil {
		fmt.Println("Step 5: Error reading:", err.Error())
		//break
	}
}

// Sents the file given through the channel given
func sendFile(connection net.Conn, fileName string, channelName string) {

	channelSendFile := ""
	if channelName == "" {
		channelSendFile = "ALL subscribed channels"
	}
	fmt.Printf("Sending %s to all clients subscribed to the following channel: %s\n", fileName, channelSendFile)
	listChannels()

	// Check if filename exceeds 64 bytes

	// Step 1: Send command
	// The protocol number for sending a file is 24
	n, err := connection.Write([]byte{24})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Sending %d byte (command)\n", n)

	//fmt.Printf("Step 1: sent %d bytes\n", n)

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
	//fmt.Printf("Step 2: the size in bytes of the file names length is %d bytes\n", fileNameInBytesSize)

	fileNameBufferLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(fileNameBufferLength, uint32(fileNameInBytesSize))
	//fmt.Printf("The buffer of the length is %v \n", fileNameBufferLength)

	// Step 2: Send file name size
	n, err = connection.Write(fileNameBufferLength)

	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	//fmt.Printf("Step 2: sent %d bytes\n", n)

	// Step 3: Send file name buffer
	n, err = connection.Write([]byte(fileName))
	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	//fmt.Printf("Step 3: sent %d bytes\n", n)

	// Step 4: Send file buffer size
	//fmt.Println("file size in bytes is", fileSize)

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

	////////////////
	// Step 6: Send channel Name length  through 4 bytes (1 int)
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

	//fmt.Printf("Step 6: sent channel name length %d bytes\n", n)

	//Step 7: Send channel name
	_, err = connection.Write(channelNameInBytes)

	if err != nil {
		fmt.Println("We couldn't send the message to the server")
		fmt.Println(err)
		return
	}

	fmt.Printf("Step 7: sent File to channel name %s bytes\n", channelName)
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

// Reads the response of the server regarding the status of the subscription to
// a channel and prints it
func subscribeToChannelResponse(connection net.Conn) {

	// Read int (4 bytes) to determine the length of the response (string)
	responseLengthBuffer := make([]byte, 4)

	// Read response length
	bytesRead, err := io.ReadFull(io.LimitReader(connection, 4), responseLengthBuffer)
	if err != nil {
		fmt.Println("Step 1 error:", err.Error())
		//break
	}
	fmt.Println(bytesRead)

	// Convert response length buffer to the length of the response
	responseLength := int32(binary.LittleEndian.Uint32(responseLengthBuffer))

	// Buffer that holds the response
	responseBuffer := make([]byte, int64(responseLength))

	// Read response
	bytesRead, err = io.ReadFull(io.LimitReader(connection, int64(responseLength)), responseBuffer)
	if err != nil {
		fmt.Println("Step 2 error:", err.Error())
		//break
	}

	response := string(responseBuffer)

	// Convert FileName buffer to string
	fmt.Printf("%s\n", response)

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

// Prints the current subscribed channels
func listChannels() {
	fmt.Printf("Channels subscribed to: %v \n", currentChannels)
}

func main() {

	// Connect with the server
	serverConnection, err := net.Dial("tcp", "0.0.0.0:8085")
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
