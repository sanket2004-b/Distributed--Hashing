package methods

import (
	"bytes"
	"distributed-hashing/client/utils/hashring"
	"distributed-hashing/client/utils/logger"
	"encoding/json"
	"fmt"
	"io"
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

func GetKeyValue(key string) ([]byte, error) {
	node := ring.GetNode(key)
	if key == "" {
		LOG.Info("key is not find")
		fmt.Printf("empty key")
	}

	if node == "" {
		err := fmt.Errorf("Unable to get node for key %s", key)
		LOG.Error(err, "Failed in getting node for key")
		return nil, err
	}
	nodeUrl := NodeToUrlMaps[node] + "/get?key=" + key

	LOG.Info("calling to thr url ::--", nodeUrl)

	resp, err := http.Get(nodeUrl)

	if err != nil {
		LOG.Error(err, "Error while calling", "url", nodeUrl)

		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		err := fmt.Errorf("key not found %s", key)
		LOG.Error(err, "Failed to get the key")
		return nil, err
	} else if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			LOG.Error(err, "Faliled to the get the data from  the response body")

		}
		LOG.Info("JSON body is ::  ", "key is:  ", key, "  body : ", string(body))

		return body, nil
	} else {
		err := fmt.Errorf("Key %v not found, invalid statuscode: %v", key, resp.StatusCode)
		return nil, err
	}

}

func DeleteKey(key string) error {
	node := ring.GetNode(key)

	if node == "" {
		err := fmt.Errorf("Unable to get node for %s", key)
		LOG.Error(err, "Unable to geting node for key")
		return err
	}
	nodeUrl := NodeToUrlMaps[node] + "/delete?key=" + key

	LOG.Info("calling  ", "url", nodeUrl)

	req, err := http.NewRequest(http.MethodDelete, nodeUrl, nil)

	if err != nil {
		LOG.Error(err, "error while req DELETE is Created ", nodeUrl)
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		LOG.Error(err, "Error while calling delete req for ", "key", key)
		return err
	}

	LOG.Info("Response ", "status", resp.Status, "key", key)
	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		LOG.Error(err, "Error while deleting ", "key", key)
		return err
	}

	LOG.Info("Response ", "Body", string(respBody))

	return nil
}
