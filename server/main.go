package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// // Client
type Client struct {
	connection net.Conn
	username   string
	server     *ServerHub
}

// Creates a new Client and returns a pointer to it
//
// Input:
// connection (net.Conn): connection with the real client
//
// Output:
// pointer to the new Client
func newClient(connection net.Conn, name string, server *ServerHub) *Client {
	return &Client{connection: connection, username: name, server: server}
}

func (client *Client) read() {
	defer client.connection.Close()

	for {

	}
}

func (client *Client) receiveFile() {
	// The protocol number for receiving a file is 24

	// Buffer that holds command protocol number
	/* commandProtocolBuffer := make([]byte, 1) */

	// Buffer that holds file name length
	fileNameLengthBuffer := make([]byte, 4)

	// Buffer that holds file  length
	fileLengthBuffer := make([]byte, 8)

	//for {
	/* // Step 1: Read command buffer
	fmt.Println("Step 1: Read command")
	bytesRead, err := io.ReadFull(io.LimitReader(client.connection, 1), commandProtocolBuffer)

	if err != nil {
		fmt.Println("Step 1 error:", err.Error())
		break
	}

	// Convert Command Buffer
	commandNumber := int(commandProtocolBuffer[0])

	fmt.Printf("Step 1: Received command: %d \n", commandNumber)

	fmt.Printf("Bytes read: %v \n", bytesRead)
	fmt.Printf("Buffer received: %v\n", commandProtocolBuffer)
	fmt.Printf("Buffer received: %v\n", commandNumber) */

	println("############# SERVER: START READ FILE #############")

	// Step 2: Read name file length
	fmt.Println("Step 2: Read name file length")

	// Read File Name length buffer
	// bytesRead, err = io.ReadFull(client.connection, fileNameLengthBuffer)
	bytesRead, err := io.ReadFull(io.LimitReader(client.connection, 4), fileNameLengthBuffer)
	if err != nil {
		fmt.Println("Step 2 error:", err.Error())
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

	bytesRead, err = io.ReadFull(io.LimitReader(client.connection, int64(fileNameLength)), fileNameBuffer)
	if err != nil {
		fmt.Println("Step 3 error:", err.Error())
		//break
	}

	// Convert FileName buffer to string
	fileName := string(fileNameBuffer)
	fmt.Printf("Step 3: The name of the file is: %s\n", fileName)

	// Step 4: Get the buffer size of the file
	bytesRead, err = io.ReadFull(io.LimitReader(client.connection, 8), fileLengthBuffer)
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
	// bytesRead, err = io.ReadFull(io.LimitReader(client.connection, fileLength), receivedFile)
	//bytes, err := io.Copy(client.connection, receivedFile)
	bytes, err := io.CopyN(receivedFile, client.connection, fileLength)
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
	println("############# SERVER: END READ FILE #############")

	client.server.sendFileToAllChannels(fileName, client)

	//}
}

func (sever *ServerHub) sendFileToAllChannels(fileName string, sender *Client) {
	//channels := make([]string, 0)
	for currentChannelName, membersSlice := range sender.server.channels {
		for _, client := range membersSlice {
			if client == sender {
				client.server.receiveAndReSendFile(fileName, currentChannelName, sender)
				break
			}
		}
	}
}

func (server *ServerHub) receiveAndReSendFile(fileName string, channelName string, sender *Client) {

	println("############# SERVER: START WRITE FILE #############")
	// Change client.server.clients to client.server.clients in a channel
	fmt.Println(server.channels)
	fmt.Println(channelName)
	fmt.Println(server.channels[channelName])
	for index, currentClient := range server.channels[channelName] {
		// if currentClient == sender {
		// 	continue
		// }
		println("@@@@@@@@Now with client", index, "@@@@@@@on@channel@@", channelName)
		fmt.Println("Current cliennts ##################:", server.clients, currentClient)
		// Check if filename exceeds 64 bytes

		// Step 1: Send command
		// The protocol number for sending a file is 24
		n, err := currentClient.connection.Write([]byte{24})
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
		n, err = currentClient.connection.Write(fileNameBufferLength)

		if err != nil {
			fmt.Println("We couldn't send the message to the server")
			fmt.Println(err)
			return
		}

		fmt.Printf("Step 2: sent %d bytes\n", n)

		// Step 3: Send file name buffer
		n, err = currentClient.connection.Write([]byte(fileName))
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

		n, err = currentClient.connection.Write([]byte(fileSizeBuffer))
		if err != nil {
			fmt.Println("We couldn't send the message to the server")
			fmt.Println(err)
			return
		}

		// Step 5: Send file buffer

		bytesWritten, err := io.Copy(currentClient.connection, file)
		if err != nil {
			fmt.Println("We couldn't send the message to the server")
			fmt.Println(err)
			return
		}
		fmt.Printf("Sent %d bytes of the file named %s \n", bytesWritten, fileName)
	}
	println("############# SERVER: END WRITE FILE #############")
}

// TODO change slice of clients in a channel for a set
func (client *Client) handleChannelSubscription() {
	// Receive channel name length
	channelNameLengthBuffer := make([]byte, 4)

	bytesRead, err := io.ReadFull(io.LimitReader(client.connection, 4), channelNameLengthBuffer)
	if err != nil {
		fmt.Println("Error receiving the length of the channel name:", err.Error())
		return
	}
	fmt.Println(bytesRead)

	channelNameLength := int32(binary.LittleEndian.Uint32(channelNameLengthBuffer))

	// Receive channel name
	channelNameBuffer := make([]byte, channelNameLength)
	bytesRead, err = io.ReadFull(io.LimitReader(client.connection, int64(channelNameLength)), channelNameBuffer)
	if err != nil {
		fmt.Println("Error receiving channel name:", err.Error())
		return
	}
	channelName := string(channelNameBuffer)
	println("Channel name is: ", channelName)

	// Subscribe the client to the channel
	// If the channel does not exist, create it
	client.server.addClientToChannel(client, channelName)
}

func (client *Client) changeUserName() {
	// Read username length
	clientUserNameLengthBuffer := make([]byte, 4)
	bytesRead, err := io.ReadFull(io.LimitReader(client.connection, 4), clientUserNameLengthBuffer)
	if err != nil || bytesRead != 4 {
		fmt.Println("Error receiving channel name:", err.Error())
		return
	}

	clientUserNameLength := int32(binary.LittleEndian.Uint32(clientUserNameLengthBuffer))

	// Read username
	clientUserNameBuffer := make([]byte, clientUserNameLength)
	bytesRead, err = io.ReadFull(io.LimitReader(client.connection, int64(clientUserNameLength)), clientUserNameBuffer)

	clientUserName := string(clientUserNameBuffer)

	// Set new username
	client.username = clientUserName
	fmt.Printf("Client new username is: %s", client.username)
}

// Handles user request by identifying if the user wants to send a message, file
// or subscribe to a certain channel. If a user subscribes to a channel that
// does not exist, it will create one
func (client *Client) handleClientRequest() {
	// Buffer that holds command protocol number
	commandProtocolBuffer := make([]byte, 1)

	for {
		// Read a single byte, this byte determines the command (request) that
		// is made to this particular client
		// If that byte is not identified as a proper command, all the
		// connection will be flushed leaving it empty. This will take care of
		// errors and unexpected behaviour

		// Step 1: Read command buffer
		fmt.Println("Step 1: Read command")
		_, err := io.ReadFull(io.LimitReader(client.connection, 1), commandProtocolBuffer)

		if err != nil {
			fmt.Println("Step 1 error:", err.Error())
			break
		}

		// Convert Command Buffer
		commandNumber := int(commandProtocolBuffer[0])

		switch commandNumber {
		// Send File
		case 24:
			client.receiveFile()
		case 34:
			client.handleChannelSubscription()
		case 44:
			client.changeUserName()
		}
	}
}

// // ServerHub
// It is responsible to handle all the connections and manage messages between
// clients and channels
type ServerHub struct {
	// Each channel will have a slice of pointers to clients
	channels map[string][]*Client

	// Store a pointer to each client
	clients []*Client

	//newClient

	//commandsChannel chan
}

// Creates a new ServerHub and returns a pointer to it
//
// Input:
// connection (net.Conn): connection with the real client
//
// Output:
// pointer to the new Client
func newServerHub() *ServerHub {

	return &ServerHub{channels: make(map[string][]*Client), clients: []*Client{}}
}

// Runs the serverHub
func (server *ServerHub) run() {

}

// When a file is sent to through the server, the server must resend this file
// to all the clients in a given channel
func (server *ServerHub) sendFile() {

}

// Add client to the server
func (server *ServerHub) addClient(client *Client) {
	server.clients = append(server.clients, client)
}

// List current clients
func (server *ServerHub) listClients() {
	fmt.Println(server.clients)
}

func (server *ServerHub) addClientToChannel(client *Client, channelName string) {
	fmt.Printf("Adding client: %v to channel: %s \n", client, channelName)
	server.channels[channelName] = append(server.channels[channelName], client)
	println("Client added!")
}

// Send buffer to all clients in the same channel
// func (server *ServerHub) sendData(channelName string, ) {

// }

func main() {

	// Create TCP server
	serverConnection, error := net.Listen("tcp", ":8080")

	// Check if an error occured
	// Note: because 'go' forces you to use each variable you declare, error
	// checking is not optional, and maybe that's good
	if error != nil {
		fmt.Println(error)
		return
	}

	// Close the server just before the program ends
	defer serverConnection.Close()

	// Create server Hub
	serverHb := newServerHub()

	// TODO this has to listen for any request made by a client, so use go so that it is always listening!
	go serverHb.run()

	// Each client sends data, that data is received in the server by a client struct
	// the client struct then sends the data, which is a request to a 'go' channel, which is similar to a queue

	// Somehow this for loop runs only when a new connection is detected
	for {

		// Accept a new connection if a request is made
		// serverConnection.Accept() blocks the for loop
		// until a connection is accepted, then it blocks the for loop again!
		connection, connectionError := serverConnection.Accept()

		// Check if an error occurred
		if connectionError != nil {
			fmt.Println("1: Woah, there's a mistake here :/")
			fmt.Println(connectionError)
			fmt.Println("1: Woah, there's a mistake here :/")
			// return
		}

		// Create new user
		var client *Client = newClient(connection, "", serverHb)
		fmt.Println(client)

		// Add client to serverHub
		serverHb.addClient(client)
		serverHb.listClients()

		// go client.receiveFile()
		go client.handleClientRequest()

	}
}
