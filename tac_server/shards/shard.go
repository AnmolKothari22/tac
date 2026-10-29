package shards

import (

	"sync"

)

func Shard_num()uint64{
	const number_of_shard=64;
	return number_of_shard;
}

type Shard struct{
	Rmu sync.RWMutex;
	Store map[string]string ;
};