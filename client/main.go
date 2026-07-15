/*store server -> url maps.

add one by one node from server in hashring.
give few set key-value calls.
also fetch those values.

set-key value function:
	convert key into hash.
	call hashring get node function.
	send key on given node.
get key value function:
	same as above.
delete key same as above.*/

package main

import (
	"distributed-hashing/client/handler"
	"distributed-hashing/client/methods"
	"os"
)

func main() {
	methods.SetUp()
	// test.UnitTesting()

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	handler.CreateHandler(port)
}
