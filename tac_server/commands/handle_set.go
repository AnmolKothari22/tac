package commands

import (
	//"fmt"
	"net"
	"sync"

	"github.com/AnmolKothari22/tac/tac_server/input_output"
)
func Handle_set(rmu *sync.RWMutex , parts []string,conn net.Conn,mp *map[string]string){
		rmu.Lock()	
		(*mp)[parts[1]]=parts[2]
		rmu.Unlock()
		conn.Write([]byte(inputoutput.Write_RESP_bulk("OK")));
}