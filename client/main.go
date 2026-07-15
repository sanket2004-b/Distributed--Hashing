package main

import (
	"distributed-hashing/client/handler"
	"distributed-hashing/client/methods"
	"distributed-hashing/testing/test"
	"os"
)

func main() {
	methods.SetUp()
	test.UnitTesting()

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	handler.CreateHandler(port)
}
