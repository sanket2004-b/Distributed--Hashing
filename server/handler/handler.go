package handler

import (
	"distributed-hashing/server/logger"
	"sync"
)

var LOG = logger.InitLogger("Logs/server.log")

var store sync.Map
var hm *robinhood.HashMap

func init() {
	hm = robinhood.Create(0.75, 16)

}

var pool *WorkerPool

func 
