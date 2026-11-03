package commands
import (
	"sync"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_del(rmu *sync.RWMutex , parts []string,mp *map[string]string) (string,inputoutput.WriteType){
	rmu.Lock()
	defer rmu.Unlock()
	//log.Print("req is processing [del]")
	var deleted_val,found=(*mp)[parts[1]]
	if !found{
		mess:="key not found"
		return mess,inputoutput.KeyMiss
	}
	delete((*mp),parts[1])
	mess:="deleted key: "+parts[1]+" with value: "+deleted_val
	return mess,inputoutput.Normal
}