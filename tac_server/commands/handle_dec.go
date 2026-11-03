package commands
import (
	"sync"
	"net"
	"strconv"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_dec(rmu *sync.RWMutex , parts []string,conn net.Conn,mp *map[string]string){
	rmu.Lock()
	defer rmu.Unlock()
	//checking key correctness
	val , found := (*mp)[parts[1]];
	//if key not found
	if !found{
		mess:="key not found"
		conn.Write([]byte(inputoutput.Write_RESP_ERR(mess)))
		return ;
	}
	//if value is NAN
	conv_num , err := strconv.ParseInt(val,10,64)
	if err!=nil{
		mess:="value is NAN"
		conn.Write([]byte(inputoutput.Write_RESP_ERR(mess)))
		return ;
	}
	var dec_by string
	if len(parts)<3{
		conv_num=conv_num-1
		dec_by="1"
	}else{
		sub_num,err := strconv.ParseInt(parts[2],10,64)
		//checking value's correctness
		if err!=nil{
			mess:="decreasing by a NAN not allowed"
			conn.Write([]byte(inputoutput.Write_RESP_ERR(mess)))
			return ;
		}
		conv_num=conv_num- sub_num
		dec_by=parts[2]
	}
	(*mp)[parts[1]]= strconv.FormatInt(conv_num,10); 
	mess:="key: "+parts[1]+" decreased by "+dec_by
	conn.Write([]byte(inputoutput.Write_RESP_bulk(mess)))
}