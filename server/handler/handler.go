package handler

import (
	"distributed-hashing/server/hashmap/robinhood"
	"distributed-hashing/server/logger"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var LOG = logger.InitLogger("Logs/server.log")

var store sync.Map
var hm *robinhood.HashMap

func init() {
	hm = robinhood.CreateHashMap(0.75, 16)

}

var pool *WorkerPool

func CreateWorkerPool(numWorkers int, hm *robinhood.HashMap) *WorkerPool {
	pool = &WorkerPool{
		Tasks:     make(chan Task),
		WorkerCnt: numWorkers,
		HM:        hm,
	}
	for i := 0; i < numWorkers; i++ {
		go pool.worker(i)
	}

	return pool
}

func (pool *WorkerPool) worker(id int) {
	for task := range pool.Tasks {
		switch task.Operation {
		case "SET":
			err := pool.HM.Put(task.Key, task.Value)
			if err != nil {
				task.Err <- fmt.Errorf("error occurred while setting key %s: %v", task.Key, err)
			} else {
				task.Result <- fmt.Sprintf("Key %s set successfully", task.Key)
			}
		default:
			task.Err <- fmt.Errorf("wrong menthod called %s", task.Operation)

		}
	}
}

func (pool *WorkerPool) AddTask(task Task) {
	pool.Tasks <- task
}

func CreateHandler(port string) {
	pool = CreateWorkerPool(20, hm)
	// fmt.Printf("\nListening on address: :%v\n", port)

	http.HandleFunc("/set", handleSet)
	address := ":" + port
	fmt.Printf("Listening on address: %v", address)
	LOG.Info("Listening on", "address", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	var req keyValRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || req.Key == "" {
		LOG.Error(err, "Invalid json or missing data", req.Key)
		http.Error(w, "Invalid json or missing data", http.StatusBadRequest)
		return
	}

	LOG.Info("Received request to set key: ", req.Key, " with value: ", req.Value)
	task := Task{
		Operation: "SET",
		Key:       req.Key,
		Value:     req.Value,
		Result:    make(chan interface{}),
		Err:       make(chan error),
	}
	LOG.Info("Adding task to worker pool for key: ", req.Key)
	pool.AddTask(task)

	select {
	case <-task.Result:
		LOG.Info("Key ", req.Key, " set successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Key %s set successfully", req.Key)))

	case err := <-task.Err:
		LOG.Error(err, "Error occurred while setting key: ", req.Key)
		http.Error(w, fmt.Sprintf("Error occurred while setting key %s: %v", req.Key, err), http.StatusInternalServerError)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		LOG.Error(nil, "missing key", "key", key)
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	LOG.Info("Fetch in", "key", key)

	task := Task{
		Operation: "GET",
		Key:       key,
		Result:    make(chan interface{}),
		Err:       make(chan error),
	}

	LOG.Info(" for fetching ", "key", key)
	pool.AddTask(task)

	select {
	case data := <-task.Result:
		val := data.([]byte)

		if !json.Valid(val) {
			http.Error(w, fmt.Sprintf("Invalid JSON stored for key: %v", key), http.StatusInternalServerError)
			return
		}
		LOG.Info("successfully got value from hashmap for ", "key", key)
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		w.Write(val)
	case err := <-task.Err:

		LOG.Error(err, "Error while fetching", "key ", key)
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
