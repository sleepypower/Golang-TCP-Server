package main

import (
	"fmt"
	"net"
)

//// Client
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

//
func (client *Client) read() {
	defer client.connection.Close()
	for {

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

// Handles user request by identifying if the user wants to send a message, file
// or subscribe to a certain channel. If a user subscribes to a channel that
// does not exist, it will create one
func (client *Client) handleClientRequest() {

}

//// ServerHub
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

		go client.read()

	}
}
