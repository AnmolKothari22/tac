package commands
import (
	"sync"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_get(rmu *sync.RWMutex , parts []string,mp *map[string]string)(string,inputoutput.WriteType){
	rmu.RLock()
	defer rmu.RUnlock()
	value , ok := (*mp)[parts[1]]
	var mess string = value;
	if ok {
    	return mess,inputoutput.Normal
	} else {
    	return mess,inputoutput.KeyMiss
	}
}