package commands
import (
	"sync"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
)
func Handle_get(rmu *sync.RWMutex , parts []string, mp *map[string]shards.Stored_values)(string,inputoutput.WriteType){
	rmu.RLock()
	defer rmu.RUnlock()

	//getting value
	value , ok := (*mp)[parts[1]]
	var mess string = value.Value;

	//if found 
	if ok {
		if(value.IsExpired()){ // if key is expired , shadow delete
			delete((*mp),parts[1])
			return "value expired", inputoutput.KeyMiss
		}
    	return mess,inputoutput.Normal // if key isnt expired, then send data
	} else {
    	return "",inputoutput.KeyMiss 
	}
}