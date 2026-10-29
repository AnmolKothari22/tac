package commands
import (
	"sync"
	"net"
	"strconv"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_inc(rmu *sync.RWMutex , parts []string,conn net.Conn,mp *map[string]string){
	rmu.Lock()
	conv_num , _ := strconv.ParseInt((*mp)[parts[1]],10,64)
	var inc_by string
	if len(parts)<3{
		conv_num=conv_num+1
		inc_by="1"
	}else{
		add_num,_ := strconv.ParseInt(parts[2],10,64)
		conv_num=conv_num+ add_num
		inc_by=parts[2]
	}
	(*mp)[parts[1]]= strconv.FormatInt(conv_num,10);
	rmu.Unlock() 
	mess:="key: "+parts[1]+" increased by "+inc_by
	conn.Write([]byte(inputoutput.Write_RESP_bulk(mess)))
}