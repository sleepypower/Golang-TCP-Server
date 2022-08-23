package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

// // Client
type Client struct {
	connection net.Conn
	username   string
	server     *ServerHub
	channels   []string
}

// Creates a new Client and returns a pointer to it
//
// Input:
// connection (net.Conn): connection with the real client
//
// Output:
// pointer to the new Client
func newClient(connection net.Conn, name string, server *ServerHub) *Client {
	return &Client{connection: connection, username: name, server: server, channels: make([]string, 0)}
}

func (client *Client) handleConnectionError(err error) {

	if err != nil {
		client.server.deleteClient(client)
	}

}

func (client *Client) receiveFile() {
	// The protocol number for receiving a file is 24
	// Buffer that holds file name length
	fileNameLengthBuffer := make([]byte, 4)

	// Buffer that holds file  length
	fileLengthBuffer := make([]byte, 8)

	println("############# SERVER: START READ FILE #############")

	// Step 2: Read name file length
	fmt.Println("Step 2: Read name file length")

	// Read File Name length buffer
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

	// Step 6
	// Read the channel name length through 4 bytes (1 int)

	// Receive channel name length
	channelNameLengthBuffer := make([]byte, 4)

	bytesRead, err = io.ReadFull(io.LimitReader(client.connection, 4), channelNameLengthBuffer)
	if err != nil {
		fmt.Println("Error receiving the length of the channel name:", err.Error())
		return
	}
	fmt.Println(bytesRead)

	channelNameLength := int32(binary.LittleEndian.Uint32(channelNameLengthBuffer))

	// Step 7
	// Receive channel name
	channelNameBuffer := make([]byte, channelNameLength)
	bytesRead, err = io.ReadFull(io.LimitReader(client.connection, int64(channelNameLength)), channelNameBuffer)
	if err != nil {
		fmt.Println("Error receiving channel name:", err.Error())
		return
	}
	channelName := string(channelNameBuffer)
	println("Azula:: Channel name is: ", channelName)

	////////
	client.server.sendFileToChannel(fileName, client, channelName)
}

// Sends the file to all clients listening to 'channelName', if 'channelName' is
// empty, sends the file to all the channels
func (sever *ServerHub) sendFileToChannel(fileName string, sender *Client, channelName string) {

	if channelName != "" {
		sender.server.receiveAndReSendFile(fileName, channelName, sender)
		fmt.Printf("Send file through channel : %s\n", channelName)
	} else {

		for currentChannelName, membersSlice := range sender.server.channels {
			for _, client := range membersSlice {
				if client == sender {
					client.server.receiveAndReSendFile(fileName, currentChannelName, sender)
					break
				}
			}
		}
		fmt.Printf("Send file through ALL channels : \n")
	}

}

func (server *ServerHub) receiveAndReSendFile(fileName string, channelName string, sender *Client) {
	fmt.Println(server.channels)
	fmt.Println(channelName)
	fmt.Println(server.channels[channelName])
	for index, currentClient := range server.channels[channelName] {

		if currentClient == sender {
			continue
		}

		println("Sending ", fileName, " to client ", index, " on channel ", channelName)
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
		currentClient.server.bytesSent += bytesWritten
		currentClient.server.filesSent += 1
		fmt.Printf("Sent %d bytes of the file named %s \n", bytesWritten, fileName)
	}
	println("############# SERVER: END WRITE FILE #############")
}

func (server *ServerHub) deleteClient(clientToBeDeleted *Client) {

	// Delete client in clients slice
	for i, currentClient := range server.clients {
		if currentClient == clientToBeDeleted {
			server.clients = removeHelper(server.clients, i)
		}
	}

	// Delete client in all the channels that it appears on
	for channelName, channel := range server.channels {
		for i, currentClient := range channel {
			if currentClient == clientToBeDeleted {
				server.channels[channelName] = removeHelper(server.clients, i)
			}
		}
	}
	clientToBeDeleted.connection.Close()
	clientToBeDeleted.connection = nil

	fmt.Printf("Client %s has been deleted!\n", clientToBeDeleted.username)

}

func removeHelper(s []*Client, i int) []*Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
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
	client.channels = append(client.channels, channelName)
}

func (client *Client) changeUserName() {
	// Read username length
	clientUserNameLengthBuffer := make([]byte, 4)
	bytesRead, err := io.ReadFull(io.LimitReader(client.connection, 4), clientUserNameLengthBuffer)
	if err != nil || bytesRead != 4 {
		fmt.Println("ChangeUserNameError:", err.Error())
		return
	}

	clientUserNameLength := int32(binary.LittleEndian.Uint32(clientUserNameLengthBuffer))

	// Read username
	clientUserNameBuffer := make([]byte, clientUserNameLength)
	bytesRead, err = io.ReadFull(io.LimitReader(client.connection, int64(clientUserNameLength)), clientUserNameBuffer)

	clientUserName := string(clientUserNameBuffer)

	// Set new username
	client.username = clientUserName
	fmt.Printf("A client has changed its username to: %s\n", client.username)
	client.server.listClients()
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
		fmt.Println("###Reading for requests###")
		_, err := io.ReadFull(io.LimitReader(client.connection, 1), commandProtocolBuffer)

		client.handleConnectionError(err)

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

func (client *Client) toJson() string {
	channelsSlice := make([]string, 0)
	for _, chn := range client.channels {
		channelsSlice = append(channelsSlice, `"`+chn+`"`)
	}
	jsonClient := `{"username":"` + client.username + `", "channels":[` + strings.Join(channelsSlice, ",") + `]}`
	return jsonClient
}

// // ServerHub
// It is responsible to handle all the connections and manage messages between
// clients and channels
type ServerHub struct {
	// Each channel will have a slice of pointers to clients
	channels map[string][]*Client

	// Store a pointer to each client
	clients []*Client

	bytesSent int64
	filesSent int
}

// Creates a new ServerHub and returns a pointer to it
//
// Input:
// connection (net.Conn): connection with the real client
//
// Output:
// pointer to the new Client
func newServerHub() *ServerHub {

	return &ServerHub{channels: make(map[string][]*Client), clients: []*Client{}, bytesSent: 0, filesSent: 0}
}

func (server *ServerHub) clientsToJson() string {
	clientsJson := `[`
	for i, client := range server.clients {
		clientsJson += client.toJson()
		if i != len(server.clients)-1 {
			clientsJson += `,`
		}
	}
	clientsJson += `]`

	return clientsJson
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
	fmt.Println("###Current Clients###")
	for _, currentClient := range server.clients {
		fmt.Printf(" - %s \n", currentClient.username)
	}
	fmt.Println("#####################")
	fmt.Printf("The number of active clients is %d", len(server.clients))
	fmt.Printf("The number of channels is %d", len(server.channels))

}

func (server *ServerHub) addClientToChannel(client *Client, channelName string) {
	fmt.Printf("Adding client: %v to channel: %s\n", client.username, channelName)
	server.channels[channelName] = append(server.channels[channelName], client)
	fmt.Println("###Client added###")
}

// Returns the number of active clients
func (server *ServerHub) getNumberOfClients() int {
	currentUsers := 0

	currentUsers = len(server.clients)

	return currentUsers

}

// Returns the number of channels
func (server *ServerHub) getNumberOfChannels() int {
	currentChannels := 0

	currentChannels = len(server.channels)

	return currentChannels

}

// Returns the number of Bytes sent
func (server *ServerHub) getBytesSent() int64 {
	var BytesSent int64
	BytesSent = 0

	BytesSent = int64(server.bytesSent)

	return int64(BytesSent)

}

// Returns the number of Files sent
func (server *ServerHub) getFilesSent() int {
	FilesSent := 0

	FilesSent = server.filesSent

	return FilesSent

}

// Get channel names for a client

func requestHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("#########\n")
	fmt.Printf("%s", serverHb.clientsToJson())
	fmt.Printf("\n#########\n")
	fmt.Println("@@@@@@@@@", serverHb.clientsToJson(), "@@@@@@@@")
	_, _ = fmt.Fprintf(w, `{ "users_connected": "%d", 
							 "files_sent": "%d", 
							 "bytes_sent": "%d", 
							 "channels": "%d",
							 "clients": %s
							}`,
		serverHb.getNumberOfClients(),
		serverHb.getFilesSent(),
		serverHb.getBytesSent(),
		serverHb.getNumberOfChannels(),
		serverHb.clientsToJson(),
	)

}

var serverHb = newServerHub()

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

	// Create server Hub

	// Close the server just before the program ends
	defer serverConnection.Close()

	// Handle Front End requests
	http.HandleFunc("/api/thumbnail", requestHandler)

	fs := http.FileServer(http.Dir("../../tcp-server-frontend/dist"))
	http.Handle("/", fs)

	fmt.Println("Server listening on port 3000")
	go http.ListenAndServe(":3000", nil)

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
		var client *Client = newClient(connection, "Unregistered_User", serverHb)
		fmt.Println(client)

		// Add client to serverHub
		serverHb.addClient(client)
		serverHb.listClients()

		// go client.receiveFile()
		go client.handleClientRequest()

	}
}
