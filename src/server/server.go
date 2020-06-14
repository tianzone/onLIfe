// Copyright 2020, Tianz. All rights reserved.

package server

import(
	"net"
)

// The constant value set
const(
	// The maximum connection is allowed for the server.
	MAX_CONNECTION = 1024
)

// A net server that responsible for accept the new connection of client.
// And despatch the message to the destination connection.
type NetServer struct{
	// The active user map
	activeUserMap map[userId]userInfo

	// The lock is used when manipulate activeUserMap
	lock sync.Mutex

	// The id generator of active users.
	userIdGenerator idGenerator

	// The id generator of conversations
	convstIdGenerator idGenerator 
}

// Set up the listening socket.
func (ns *NetServer) StartupAndServe( addr, port string ) error{
	// TODO: Process the address and port when any of them is invalid
	finalAddr = addr + ":" + port
	fmt.Printf( "The address of listening: %s\n", finalAddr )

	// Set up the listening socket.
	listener, err : = net.Listen( "tcp", finalAddr )
	if err != nil{
		fmt.Printf( "Server attempts to listen %s fialed: %s\n", finalAddr, err )
		return err
	}

	// A loop to accept the new connection of client
	for{
		newConn, err := listener.Accept()
		if err != nil{
			fmt.Printf( "Attempts to accept the new connection failed: %s\n", err )
			// TODO: Proccess the error with the specific case. Terminate the program if
			// it's the fatal errors. Otherwise continue to accept next one connection.
			continue;
		}

		// Add the new connetion into the active connection list.
		go addNewConnection( newConn )
	}

	return nil
}

// Generates the uniqe id to the specifies connection. Then map them on the active user map.
func (ns *NetServer) addNewConnection( conn net.Conn ){
	// Print out the remote host infos
	fmt.Printf( "Add new connection: addr: %s\n", conn.RemoteAddr().String() );

	newUserId, err := ns.userIdGenerator.newId()
	if err != nil{
		return;
	}else{
		ns.lock.Lock()
		defer ns.lock.Unlock()
		ns.activeUserMap[newUserId] := userInfo{
			id: newUserId,
			addr: conn.RemoteAddr(),
			conn: conn
		}
	}

	// Handles the conversation with the remote connection
	go handleWithRemoteConvst( &ns.activeUserMap[newUserId] )
}

// Remove the connection with the specifies id
func (ns *NetServer) removeConnection( id userId ){
	// Check to see if the user id existed.
	ns.lock.Lock()
	defer ns.lock.Unlock()
	if _, ok := ns.activeUserMap[id]; ok != nil{
		return
	}else{
		// Release the user id before remove it.
		ns.userIdGenerator.release( id )
		// Remove the user
		delete( ns.activeUserMap, id )
	}
}

// Handles the conversation with the remote connection
func (ns *NetServer) handleWithRemoteConvst( pUser *userInfo ){
	// Block to wait the remote message comming.
	for{
		// The message buffer on byte
		buff := make( []byte, 128 )
		// Reads the message from remote peer.
		if readBytes, err := pUser.conn.Read( buff ); err != nil{
			// TODO: Handles the error cases
			fmt.Println( "Read connection failed: ", error )
			continue
		}

		// Parse the conversation message
		msg := msgPac( buff )
		switch msg.opType{
			// Handles the disconnect request
			case OP_TYPE_DISCONNECT:{
				ns.removeConnection( pUser.id )
			}
			// It means the sender ask to open a new conversation.
			case OP_TYPE_OPEN_CONVST:{
				// Get a new conversation id first.
				if cId, err := ns.convstIdGenerator.newId(); err != nil{
					fmt.Println( "Can not open new conversation: ", error )
					continue
				}else{
					// Find out the target user.
					// TODO: Get the target user id from the payload data.
					uId := userId( msg.data )
					if _, isExisted := ns.activeUserMap[uId]; isExisted == true{
						// Add the conn of target user to this user's conversation list
						pUser.AddConvst( cId, &ns.activeUserMap[uId].conn )
						// Add the conn of this user to the target user's conversation list
						ns.activeUserMap[uId].AddConvst( cId, &pUser.conn )
					}
				}
			}
			// It means the sender ask to close the new conversation.
			case OP_TYPE_CLOSE_CONVST:{
				ns.removeConnection( pUser.id )
			}
			// Route the message to the peer which map with the conversation id
			case OP_TYPE_SEND_MSG:{
				ns.removeConnection( pUser.id )
			}
		}
	}
}
