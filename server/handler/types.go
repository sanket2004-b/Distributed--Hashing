package handler

import (
	"distributed-hashing/server/hashmap/robinhood"
	"sync"
)

type keyValRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Task struct {
	ID        string           `json:"id"`
	Key       string           `json:"key"`
	Value     interface{}      `json:"value"`
	Operation string           `json:"operation"`
	Result    chan interface{} `json:"result"`
	Err       chan error       `json:"err"`
}

type WorkerPool struct {
	Tasks     chan Task
	WorkerCnt int
	HM        *robinhood.HashMap
	wg        sync.WaitGroup
}
