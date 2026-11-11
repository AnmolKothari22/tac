package commands

import (
	//"fmt"
	"sync"
	"time"
	"strconv"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
)
func Handle_set(rmu *sync.RWMutex , parts []string,mp *map[string]shards.Stored_values)(string,inputoutput.WriteType){
		rmu.Lock()	
		defer rmu.Unlock()

		if len(parts)==3{ //non time bound set
			(*mp)[parts[1]]=shards.Stored_values{Value: parts[2] , Created_at: time.Now(),Expire_after: -1} // in shards, expire_after <0 implies non  ttl set
			return "OK",inputoutput.Normal
		}else{
			//duration extraction
			dur , err := strconv.ParseInt(parts[4],10,64)
			if err != nil{
				mess:="invalid time duration"
				return mess,inputoutput.Error
			}
			var conv_dur time.Duration // duration which is going to be converted to ms or s

			switch parts[3]{ // milliseconds or seconds , others lead to error
			case "ms":
				conv_dur =time.Duration(dur) * time.Millisecond
			case "s":
				conv_dur =time.Duration(dur) * time.Second
			default :
				mess:="use ms or s as 4th argument"
				return mess,inputoutput.Error
			}

			(*mp)[parts[1]]=shards.Stored_values{Value: parts[2] , Created_at: time.Now() , Expire_after : conv_dur }
			return "OK",inputoutput.Normal

		}
}