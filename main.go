package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

// // Client
type Client struct {
	connection net.Conn
	username   string
}

// Creates a new Client and returns a pointer to it
//
// Input:
// connection (net.Conn): connection with the real client
//
// Output:
// pointer to the new Client
func newClient(connection net.Conn, name string) *Client {
	return &Client{connection: connection, username: name}
}

func (client *Client) read() {
	defer client.connection.Close()

	for {

		fmt.Println("Initating read")
		// Read first byte to determine the command
		// io.LimitReader()
		buffer := make([]byte, 0, 2)
		n, err := io.ReadFull(client.connection, buffer)
		if err != nil {
			print("wAIT A goddamn second")
		}
		fmt.Println(n)
		fmt.Printf("%v", buffer)

		fmt.Println("Ending read")
		// Read first byte, the character of the first byte determines what command
		// was given

		// // Here I should replace 'ServerFileReceived.jpg' with the name of the
		// // file that is coming
		// file, err := os.Create("ServerFileReceived.txt")

		// if err != nil {
		// 	fmt.Println("2: Woah, there's a mistake here :/")
		// 	fmt.Println(err)
		// }

		// // make sure to close the file
		// defer file.Close()

		// n, err := io.Copy(file, client.connection)
		// if err != nil {
		// 	fmt.Println("3: Woah, there's a mistake here :/")
		// 	fmt.Println(err)
		// }

		// fmt.Println("########")
		// fmt.Println(n, "bytes received")
		// fmt.Println("########")

		// // Read Client Data
		// _, error := bufio.NewReader(client.connection).ReadString('\n')
		// if error != nil {
		// 	fmt.Println("4: Woah, there's a mistake here :/")
		// 	fmt.Println(error)
		// 	fmt.Println("4: Woah, there's a mistake here :/")

		// }

	}
}

func (client *Client) receiveFile() {
	// The protocol number for receiving a file is 24

	// Buffer that holds command protocol number
	commandProtocolBuffer := make([]byte, 1)

	// Buffer that holds file name length
	fileNameLengthBuffer := make([]byte, 4)

	// Buffer that holds file  length
	fileLengthBuffer := make([]byte, 8)

	for {
		// Step 1: Read command buffer
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
		fmt.Printf("Buffer received: %v\n", commandNumber)

		// Step 2: Read name file length
		fmt.Println("Step 2: Read name file length")

		// Read File Name length buffer
		// bytesRead, err = io.ReadFull(client.connection, fileNameLengthBuffer)
		bytesRead, err = io.ReadFull(io.LimitReader(client.connection, 4), fileNameLengthBuffer)
		if err != nil {
			fmt.Println("Step 2 error:", err.Error())
			break
		}

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
			break
		}

		// Convert FileName buffer to string
		fileName := string(fileNameBuffer)
		fmt.Printf("Step 3: The name of the file is: %s\n", fileName)

		// Step 4: Get the buffer size of the file
		bytesRead, err = io.ReadFull(io.LimitReader(client.connection, 8), fileLengthBuffer)
		if err != nil {
			fmt.Println("Step 4: Error reading:", err.Error())
			break
		}

		// Convert the buffer size slice to a number
		fileLength := int64(binary.LittleEndian.Uint64(fileLengthBuffer))

		fmt.Printf("Step 4: We should receive a file of size: %d bytes\n", fileLength)

		// Step 5
		// Read the buffer and copy it to the created file with name fileName
		// Create the file
		receivedFile, err := os.Create("Server_" + fileName)
		defer receivedFile.Close()

		fmt.Println("Copying...")

		// Read the file and copy it into fileName
		// bytesRead, err = io.ReadFull(io.LimitReader(client.connection, fileLength), receivedFile)
		//bytes, err := io.Copy(client.connection, receivedFile)
		bytes, err := io.CopyN(receivedFile, client.connection, fileLength)
		if err != nil {
			fmt.Println("Step 5: Error reading:", err.Error())
			break
		}
		fmt.Printf("Step 5: Bytes read: %d", bytes)
		fmt.Println("Copied successfully")

		if err != nil {
			fmt.Println("Step 5: Error reading:", err.Error())
			break
		}

	}
}

// Handles user request by identifying if the user wants to send a message, file
// or subscribe to a certain channel. If a user subscribes to a channel that
// does not exist, it will create one
func (client *Client) handleClientRequest() {

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

// Add client to the server
func (server *ServerHub) addClient(client *Client) {
	server.clients = append(server.clients, client)
}

// List current clients
func (server *ServerHub) listClients() {
	fmt.Println(server.clients)
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
		var client *Client = newClient(connection, "")
		fmt.Println(client)

		// Add client to serverHub
		serverHb.addClient(client)
		serverHb.listClients()

		// go client.read()
		go client.receiveFile()

	}
}
