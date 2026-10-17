package handle 

import (
	"bufio"
//	"fmt"
	"log"
	"net"
	"strconv"
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
		parts[0]=strings.ToLower(parts[0])
		if parts[0]=="set"{
			mu.Lock()
			mess:= parts[1] + " is key for "+parts[2]+"\n"
			conn.Write([]byte(mess))
			mp[parts[1]]=parts[2]
			mu.Unlock()
		}
		if parts[0]=="get"{
			rmu.RLock()
			log.Print("req is processing [get]")
		//	time.Sleep(10*time.Second)
			mess:="key: "+parts[1]+" value: "+mp[parts[1]]+"\n"
			conn.Write([]byte(mess))
			rmu.RUnlock()
		}
		if parts[0]=="del"{
			rmu.RLock()
			log.Print("req is processing [del]")
			var deleted_val=mp[parts[1]]
			delete(mp,parts[1])
			mess:="deleted key: "+parts[1]+" with value: "+deleted_val+"\n"
			conn.Write([]byte(mess))
			rmu.RUnlock()
		}
		if parts[0]=="inc"{
			mu.Lock()
			conv_num , _ := strconv.ParseInt(mp[parts[1]],10,64)
			var inc_by string
			if len(parts)<3{
				conv_num=conv_num+1
				inc_by="1"
			}else{
				add_num,_ := strconv.ParseInt(parts[2],10,64)
				conv_num=conv_num+ add_num
				inc_by=parts[2]
			}
			mp[parts[1]]= strconv.FormatInt(conv_num,10);
			mess:="key: "+parts[1]+" increased by "+inc_by+"\n"
			conn.Write([]byte(mess))
			mu.Unlock() 
		}
    }
}