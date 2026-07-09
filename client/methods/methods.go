package methods

import "distributed-hashing/client/utils/logger"

var LOG = logger.InitLogger("Logs/client.log")

var NodeToMaps = map[string]string{
	"hypervm-1": "http://localhost:8081",
	"hypervm-2": "http://localhost:8082",
	"hypervm-3": "http://localhost:8083",
}

type keyValRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
