package handle 

import (
	"bufio"
//	"fmt"
	"log"
	"net"

	"strings"
	"sync"
	
//	"time"

)

func Handle(conn net.Conn, mp map[string]string,wg *sync.WaitGroup,mu *sync.Mutex,rmu *sync.RWMutex) {
	defer conn.Close()
	defer wg.Done()
	log.Print("req received")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
        parts := strings.Fields(scanner.Text())
		if parts[0]=="set"{
			mu.Lock()
			mess:= parts[1] + " is key for "+parts[2]+"\n"
			conn.Write([]byte(mess))
			mp[parts[1]]=parts[2]
			mu.Unlock()
		}
		if parts[0]=="get"{
			rmu.RLock()
			log.Print("req is processing")
		//	time.Sleep(10*time.Second)
			mess:="key: "+parts[1]+" value: "+mp[parts[1]]+"\n"
			conn.Write([]byte(mess))
			rmu.RUnlock()
		}
    }
}