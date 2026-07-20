package handler

import (
	"distributed-hashing/client/methods"
	"distributed-hashing/client/utils/logger"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// import "distributed-hashing/server/logger"

var LOG = logger.InitLogger("Logs/client.log")
var setCmd = `curl -X POST "http://localhost:8080/set?key=user123" \
     -H "Content-Type: application/json" \
     -d '{
           "data": {
               "name": "Sanket",
               "age": 21,
               "tags": ["go", "dev"],
               "meta": {
                   "active": true,
                   "lastLogin": "2026-01-01T12:00:00Z"
               }
           }
         }'`

var getCMD = `curl.exe GET "http://localhost:8080/get?key=user123"`

func CreateHandler(port string) {
	http.HandleFunc("/set", HandleSet)
	http.HandleFunc("/get", handleGet)
	address := ":" + port
	fmt.Printf("\nListening on address: %v\n", address)

	fmt.Printf("Format to store key-value pair command\n %v\n \n", setCmd)
	fmt.Printf("Format to get the data from server using the key commnad is", getCMD)
	LOG.Info("Listening on address: %v", address)

	log.Fatal(http.ListenAndServe(address, nil))
}

func HandleSet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if key == "" {
		LOG.Error(nil, "Missing key parameter in request")
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
		return
	}

	var value interface{}

	err := json.NewDecoder(r.Body).Decode(&value)
	if err != nil {
		LOG.Error(err, "Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err = methods.SetKeyValue(key, value)
	LOG.Info("Set key-value pair for key: %s, value: %v", key, value)
	if err != nil {
		LOG.Error(err, "Error setting key-value pair: %v", err)
		http.Error(w, "Error setting key-value pair", http.StatusInternalServerError)
		fmt.Fprintf(w, "Error setting key-value pair: %v", err)
		return

	}
	LOG.Info("Successfully set key-value pair for key: %s, value: %v", key, value)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully set key-value pair for key: %s, value: %v", key, value)

}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		LOG.Error(nil, "key is missing or wrong req", "key", key)
		http.Error(w, "Missing key", http.StatusBadRequest)
		return

	}

	data, err := methods.GetKeyValue(key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	LOG.Info("successfully got value from the hashMap for ", "key is ", key)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
