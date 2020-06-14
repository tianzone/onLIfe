// Copyright 2020, Tianz. All rights reserved.

package server

// TODO: Introduces what is the module?

import(
	"fmt"
	"sync"
)

// The constant set
const(
	// The default size of id list
	ID_DEFAULT_SIZE = 12
)

// idGenerator will manage the id set
type idGenerator struct{
	// The base number that the generator increase base on.
	base uint64

	// The maximum number allows the generator to generates base on the base.
	size uint32

	// The map to record if a id has already been used.
	idMap map[uint64]bool

	// The lock is used when manipulate the idMap
	lock sync.Mutex
}

// Return a new id that does not be used. return nil of there're no valid id.
func (ig *idGenerator) newId() (uint64, error){
	// Initialize the id map at the first time enter the function
	sync.On( func(){
		MAX_OF_UINT64 := uint64( 0xFFFFFFFFFFFFFFFF )
		// Check to see if the maximun number exceed MAX_OF_UINT64
		ig.lock.Lock()
		defer ig.lock.Unlock()
		if ig.size == 0{
			ig.size = ID_DEFAULT_SIZE
		}else if (ig.base + ig.size) > MAX_OF_UINT64{
			ig.size = MAX_OF_UINT64 - ig.base
		}

		// Init to false.
		maxNum = ig.base + ig.size - 1
		for index := ig.base; index < maxNum; index++{
			ig.idMap[index] = false
		}
	})

	// Travel the id map and return first index that does not be used.
	ig.lock.Lock()
	defer ig.lock.Unlock()
	for index, isUsed := range( ig.idMap ); isUsed == false{
		// Record the id as in use.
		ig.idMap[index] = true
		return index, nil
	}

	return nil, fmt.Errorf( "There're no valid id" )
}

// Release a id
func (ig *idGenerator) release( id uint64 ){
	ig.lock.Lock()
	defer ig.lock.Unlock()
	delete( id.idMap, id )
}