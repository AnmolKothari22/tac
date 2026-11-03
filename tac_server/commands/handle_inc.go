package commands
import (
	"sync"
	"strconv"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_inc(rmu *sync.RWMutex , parts []string,mp *map[string]string)(string,inputoutput.WriteType){
	rmu.Lock()
	defer rmu.Unlock()
	//checking key correctness
	val , found := (*mp)[parts[1]];
	//if key not found
	if !found{
		mess:="key not found"
		return mess,inputoutput.KeyMiss ;
	}
	//if value is NAN
	conv_num , err := strconv.ParseInt(val,10,64)
	if err!=nil{
		mess:="value is NAN"
		return mess,inputoutput.Error;
	}
	var inc_by string
	if len(parts)<3{
		conv_num=conv_num+1
		inc_by="1"
	}else{
		add_num,err := strconv.ParseInt(parts[2],10,64)
		//checking value's correctness
		if err!=nil{
			mess:="increasing by a NAN not allowed"
			
			return mess,inputoutput.Error;
		}
		conv_num=conv_num+ add_num
		inc_by=parts[2]
	}
	(*mp)[parts[1]]= strconv.FormatInt(conv_num,10); 
	mess:="key: "+parts[1]+" increased by "+inc_by
	return mess, inputoutput.Normal
}