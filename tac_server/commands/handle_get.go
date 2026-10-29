package commands
import (
	"sync"
	"net"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_get(rmu *sync.RWMutex , parts []string,conn net.Conn,mp *map[string]string){
	rmu.RLock()
	value , ok := (*mp)[parts[1]]
	rmu.RUnlock()
	var mess string = value;
	if ok {
    	conn.Write([]byte(inputoutput.Write_RESP_bulk(mess)))
	} else {
    	conn.Write([]byte(inputoutput.MissingKey()))
	}
}