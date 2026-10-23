package commands
import (
	"sync"
	"net"
	"strconv"
)
func Handle_dec(rmu *sync.RWMutex , parts []string,conn net.Conn,mp map[string]string){
	rmu.Lock()
	conv_num , _ := strconv.ParseInt(mp[parts[1]],10,64)
	var dec_by string
	if len(parts)<3{
		conv_num=conv_num-1
		dec_by="1"
	}else{
		add_num,_ := strconv.ParseInt(parts[2],10,64)
		conv_num=conv_num- add_num
		dec_by=parts[2]
	}
	mp[parts[1]]= strconv.FormatInt(conv_num,10);
	rmu.Unlock() 
	mess:="key: "+parts[1]+" decreased by "+dec_by+"\n"
	conn.Write([]byte(mess))
}