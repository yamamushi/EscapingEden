# Eden Networking

Eden Networking is a standalone component of the Eden project that is response for handling all communication with clients.

It handles the setup of the connection, and manages the communication with the client as well as the ui.Console component.

Messages are received from the Console via a go channel, and are then processed. Messages will result in different
actions being taken, such as communication with a console for a specific client ID, or sending messages into the 
game server for processing. 

Messages can also be received from the game server for processing in the same manner. 

Connection will launch a go routine that will 