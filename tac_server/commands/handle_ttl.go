package commands

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
)

func Handle_ttl(rmu *sync.RWMutex , parts []string, mp *map[string]shards.Stored_values)(string,inputoutput.WriteType){
	rmu.RLock()
	defer rmu.RUnlock()
	value , ok := (*mp)[parts[1]]
	if ok{
		if value.Expire_after<0{
			mess:="no exp"
			return mess,inputoutput.Error
		}
		if value.IsExpired(){
			delete((*mp),parts[1])
			mess:="value missing"
			return mess,inputoutput.KeyMiss
		}

		
		time_rem := time.Until(value.Created_at.Add(value.Expire_after))

		if len(parts)==2 {
			elapsed:=fmt.Sprintf("%d ms", time_rem.Milliseconds())
			return elapsed,inputoutput.Normal
		}else{
			parts[2]=strings.ToLower(parts[2])

			switch parts[2]{
				case "s":
					elapsed:=fmt.Sprintf("%.2f s", time_rem.Seconds())
					return elapsed,inputoutput.Normal
				case "ms":
					elapsed:=fmt.Sprintf("%d ms", time_rem.Milliseconds())
					return elapsed,inputoutput.Normal
				default:
					return "3rd arg , ms or s only",inputoutput.Error
			}
		}

	}else{
		return "",inputoutput.KeyMiss
	}
}