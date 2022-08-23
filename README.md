# Golang-TCP-Server

Implementation of a TCP-SERVER using Golang

## Usage

### Run Server
```cmd
cd server
go run main.go
```
Allow access to Windows Defender Firewall
This will allow the server.go to list to localhost:8080

### Run Clients
There are 5 Clients, but you can create more!
```cmd
cd client'Number'
go run client.go
```

## Commands
There following commands are currently implemented HELP SEND SUB USERNAME CHNLS MSG
| Command        | Usage           | What does it do  |
| ------------- |:-------------:|:-----|
| SEND      |SEND 'filename.extension' [channelName] |Sends the file to all the channels that the current client is subscribed to. OPTIONAL [channelName]: If given a channel Name, the file will be sent only to the clients subcribed to that channel|
| SUB      |SUB 'channelName'|Subscribes the client to the channel given, all files sent to that channel will be received (if a certain client is subscribed to it) |
| USERNAME |USERNAME 'newUsername'|Changes the client's username to 'newUsername'|
| CHNLS |USERNAME|Lists all the channels that the current client is subscribed to|
| HELP |HELP|Shows all available commands|

There is still a command yet to be implemented:
| Command        | Usage           | What does it do  |
| ------------- |:-------------:|:-----|
| MSG |MSG 'messageBody' 'channelName'|SENDS a message to all the clients subscribed to channel 'channelName'|

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
