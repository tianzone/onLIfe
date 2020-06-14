// Copyright 2020, Tianz. All rights reserved.

package server

import(
	"sync"
	"net"
)

// Global varirable set
var(
	// The latest id of user
	currentId userId

	// The mutex lock is used when generate the user id.
	idLock sync.Mutex
)

// The infomation set of remote host that as a connetion exist on server.
type userId = uint64
type userInfo struct{
	// id is the uniqe number of per connection
	id userId

	// addr is the IP address of remote host
	addr net.Addr

	// connt is the conection entity that bewteen the remote host and server
	conn net.Conn

	// The conversation list with the user
	conversationList
}
