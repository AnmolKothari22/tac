package commands
import (
	"sync"
	"net"
)
func Handle_get(rmu *sync.RWMutex , parts []string,conn net.Conn,mp map[string]string){
	rmu.RLock()
	value , ok := mp[parts[1]]
	rmu.RUnlock()
	var mess string
	if ok{
		mess="value: "+value+"\n"
		}else{
			mess="<nil>\n"
		}
	conn.Write([]byte(mess))
}