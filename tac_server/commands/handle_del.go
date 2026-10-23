package commands
import (
	"sync"
	"net"
)
func Handle_del(rmu *sync.RWMutex , parts []string,conn net.Conn,mp map[string]string){
	rmu.Lock()
	//log.Print("req is processing [del]")
	var deleted_val=mp[parts[1]]
	delete(mp,parts[1])
	rmu.Unlock()
	mess:="deleted key: "+parts[1]+" with value: "+deleted_val+"\n"
	conn.Write([]byte(mess))
}