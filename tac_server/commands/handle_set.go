package commands
import (
	"sync"
	"net"
)
func Handle_set(rmu *sync.RWMutex , parts []string,conn net.Conn,mp map[string]string){
		rmu.Lock()	
		mp[parts[1]]=parts[2]
		rmu.Unlock()
		mess := "OK\n"
		conn.Write([]byte(mess))
}