package commands
import (
	"sync"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
)
func Handle_get(rmu *sync.RWMutex , parts []string, mp *map[string]shards.Stored_values)(string,inputoutput.WriteType){
	rmu.RLock()
	defer rmu.RUnlock()
	value , ok := (*mp)[parts[1]]
	var mess string = value.Value;
	if ok {
		if(value.IsExpired()){
			delete((*mp),parts[1])
			return "value expired", inputoutput.KeyMiss
		}
    	return mess,inputoutput.Normal
	} else {
    	return "",inputoutput.KeyMiss
	}
}