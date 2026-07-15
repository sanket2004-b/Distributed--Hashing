// For first cut will use go's inbuilt map.
//	create two handlers get and set.

// Set-> set key value in map.-> we will have to add LOCK here so that no one can read/write value from map.

// get-> return value from map.-> multiple reader can be allowed. But when write is going on we cant read.

package main

import (
	"distributed-hashing/server/handler"
	"os"
)

func main() {

	port := "8081"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	handler.CreateHandler(port)
}
