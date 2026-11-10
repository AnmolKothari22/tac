package shards

import (
	"sync"
	"time"
)

type Stored_values struct {
	Value string ;
	Created_at time.Time ;
	Expire_after time.Duration ;
};

func (s Stored_values) IsExpired() bool {
	if s.Expire_after<0{
        return false
    }
	expiry := s.Created_at.Add(s.Expire_after);
	if(time.Now().Before(expiry)){
		return false;
	}
	return true;
}



type Shard struct{
	Rmu sync.RWMutex;
	Store map[string]Stored_values;
};