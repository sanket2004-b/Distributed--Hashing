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
