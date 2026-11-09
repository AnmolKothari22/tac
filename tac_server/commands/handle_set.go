package commands

import (
	//"fmt"
	"sync"

	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_set(rmu *sync.RWMutex , parts []string,mp *map[string]string)(string,inputoutput.WriteType){
		rmu.Lock()	
		defer rmu.Unlock()
		(*mp)[parts[1]]=parts[2]
		return "OK",inputoutput.Normal
}