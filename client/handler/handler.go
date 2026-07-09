package handler

import (
	"distributed-hashing/client/utils/logger"
	"net/http"
)

// import "distributed-hashing/server/logger"

var LOG = logger.InitLogger("Logs/client.log")

func CreateHandler(port string) {
	http.HandleFunc("/set", HandleSet)
}
