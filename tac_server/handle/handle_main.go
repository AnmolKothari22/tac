package handle 

import (
	"bufio"
//	"fmt"
//	"log"
	"net"
	"strings"
	"sync"
	"github.com/AnmolKothari22/tac/tac_server/commands"
//	"time"

)

func Handle(conn net.Conn, mp map[string]string,wg *sync.WaitGroup,rmu *sync.RWMutex) {
	defer conn.Close()
	defer wg.Done()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
        parts := strings.Fields(scanner.Text())
		//log.Print(parts[0],"--",parts[1])
		parts[0]=strings.ToLower(parts[0])
		if parts[0]=="set"{
			commands.Handle_set(rmu,parts,conn,mp);
		}else if parts[0]=="get"{
			commands.Handle_get(rmu,parts,conn,mp);
		}else if parts[0]=="del"{
			commands.Handle_del(rmu,parts,conn,mp);
		}else if parts[0]=="inc"{
			commands.Handle_inc(rmu,parts,conn,mp);
		}else if parts[0]=="dec"{
			commands.Handle_dec(rmu,parts,conn,mp);
		}else{
			conn.Write([]byte("write correct command\n"))
		}
    }
}