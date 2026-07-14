package hashring

import (
	"crypto/sha256"
	"distributed-hashing/client/utils/logger"
	"sort"
	"strconv"
	"sync"
)

var virtualNodes = 10

type HashRing struct {
	nodes        map[uint64]string
	sortedhashes []uint64
	mutex        sync.RWMutex
}

var LOG = logger.InitLogger("Logs/hashring.log")

func CreateNewHashRing() *HashRing {
	hashring := &HashRing{

		nodes:        make(map[uint64]string),
		sortedhashes: make([]uint64, 0),
	}
	return hashring
}

// This function converts a string key into a 64-bit hash value using the SHA-256 cryptographic hash function
func ConvertKeyToHash(key string) uint64 {
	sum := sha256.Sum256([]byte(key))
	return uint64(sum[0])<<56 | uint64(sum[1])<<48 | uint64(sum[2])<<40 | uint64(sum[3])<<32 |
		uint64(sum[4])<<24 | uint64(sum[5])<<16 | uint64(sum[6])<<8 | uint64(sum[7])
}

// adding a node to the hash ring with virtual nodes
func (hr *HashRing) AddNode(node string) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()

	for i := 0; i < virtualNodes; i++ {
		virtualNodeKey := node + "_" + strconv.Itoa(i)
		hashKeyOfVirtualNode := ConvertKeyToHash(virtualNodeKey)
		hr.nodes[hashKeyOfVirtualNode] = node
		hr.sortedhashes = append(hr.sortedhashes, hashKeyOfVirtualNode)
	}
	sort.Slice(hr.sortedhashes, func(i, j int) bool {
		return hr.sortedhashes[i] < hr.sortedhashes[j]
	})
	LOG.Info("Node %s added to the hash ring with %d virtual nodes", node, virtualNodes)
}

func FindTargetedNodeHash(sortedhashes []uint64, hashKey uint64) uint64 {
	len := len(sortedhashes)
	if len == 1 {
		return sortedhashes[0]
	}

	start := 0
	end := len - 1

	for start <= end {
		mid := (start + end) / 2
		if sortedhashes[mid] == hashKey {
			return sortedhashes[mid]
		} else if sortedhashes[mid] < hashKey {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
	if start < len {
		return sortedhashes[start]
	}
	return sortedhashes[0]

}
func (h *HashRing) GetNode(node string) string {
	if len(h.nodes) == 0 {
		LOG.Info("No nodes in the hash ring", node)
		return ""
	}
	if len(h.sortedhashes) == 0 {
		LOG.Info("No sorted hashes in the hash ring", node)
		return ""
	}
	if node == "" {
		LOG.Info("Empty node name provided", node)
		return ""
	}
	hashKey := ConvertKeyToHash(node)
	targetedHash := FindTargetedNodeHash(h.sortedhashes, hashKey)
	LOG.Info("Node %s found for key %s", h.nodes[targetedHash], node)
	return h.nodes[targetedHash]
}
