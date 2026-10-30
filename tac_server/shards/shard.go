package shards

import (
	"sync"
)

type Shard struct{
	Rmu sync.RWMutex;
	Store map[string]string ;
};