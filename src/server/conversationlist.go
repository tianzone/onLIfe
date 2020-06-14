// Copyright 2020, Tianz. All rights reserved.

package server

// TODO: Introduces what is the module?

import(
	"fmt"
	"net"
	"sync"
)

// The constant set
const(
	// The operation type of per conversation
	// It means the sender ask to disconnect to the receiver.
	OP_TYPE_DISCONNECT = 0X01
	// It means the sender ask to open a new conversation.
	OP_TYPE_OPEN_CONVST = 0X02
	// It means the sender ask to close the new conversation.
	OP_TYPE_CLOSE_CONVST = 0X03
	// It means the sender send to message to peer.
	OP_TYPE_SEND_MSG = 0X04
)

// The message package
type opType_t = uint
type msgPac struct{
	// Conversation id
	id convstId

	// The type of conversation
	opType opType_t

	// The payload on bytes
	data []byte
}

// conversationList will stores all conversation pair that maps the
// conversation id and the net.Conn of remote peer.
type convstId uint
type convstMap = map[convstId]*net.Conn
type conversationList struct{
	// The active conversation map of per user.
	activeConvsts convstMap

	// The lock is used when manipulates the activeConvst
	lock sync.Mutex
}

// Add a new conversation
func (cl *conversationList) AddConvst( id convstId, pConn *net.Conn ) error{
	// Check to see if the connection is valid.
	if pConn == nil{
		return fmt.Errorf( "Invalid connection sepecified" )
	}

	// Check to see if the conversation has already existed.
	cl.lock.Lock()
	defer cl.lock.Unlock()
	if _, ok := cl.activeConvsts[id]; ok == nil{
		return fmt.Errorf( "The %d conversation has already existed", id )
	}

	cl.activeConvsts[id] = pConn
	return nil
}

// Remove the specifies connection
func (cl *conversationList) RemoveConvst( id convstId ){
	cl.lock.Lock()
	defer cl.lock.Unlock()
	delete( cl.activeConvsts, id )
}

// Return the remote connection corresponding with the secifies conversation if it existed.
// Return nil if it does not exist.
func (cl *conversationList) GetRemoteConn( id convstId ) (*net.Conn, error){
	cl.lock.Lock()
	defer cl.lock.Unlock()
	if _, ok := cl.activeConvsts[id]; ok != nil{
		return nil, fmt.Errorf( "The %d conversation did not exsit", id )
	}

	return cl.activeConvsts[id], nil
}