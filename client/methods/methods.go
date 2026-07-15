package methods

import (
	"bytes"
	"distributed-hashing/client/utils/hashring"
	"distributed-hashing/client/utils/logger"
	"encoding/json"
	"net/http"
)

var LOG = logger.InitLogger("Logs/client.log")

var NodeToUrlMaps = map[string]string{
	"hypervm-1": "http://localhost:8081",
	"hypervm-2": "http://localhost:8082",
	"hypervm-3": "http://localhost:8083",
}

type keyValRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var ring *hashring.HashRing

func SetUp() {
	ring = hashring.CreateNewHashRing()

	for nodeName := range NodeToUrlMaps {
		ring.AddNode(nodeName)
	}
}

func SetKeyValue(key string, value interface{}) error {
	node := ring.GetNode(key)

	nodeUrl := NodeToUrlMaps[node] + "/set"

	keyValReq := keyValRequest{
		Key:   key,
		Value: value,
	}

	data, err := json.Marshal(keyValReq)
	if err != nil {
		LOG.Error(err, "Error while marshalling key value request: %v")
		return err
	}

	LOG.Info("Sending request to node %s at URL %s with data: %s", node, nodeUrl, string(data))

	res, err := http.Post(nodeUrl, "application/json", bytes.NewBuffer(data))

	if err != nil {
		LOG.Error(err, "Error while sending request to node %s at URL %s: %v", node, nodeUrl, err)
		return err
	}
	LOG.Info("Response for store key: %s, response: %v", key, res)
	LOG.Info("Successfully sent request to node %s at URL %s", node, nodeUrl)
	res.Body.Close()
	return nil
}
