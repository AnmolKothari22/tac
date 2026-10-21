package handle 

import (
	"bufio"
//	"fmt"
//	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	
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
			rmu.Lock()
			//mess:= parts[1] + " is key for "+parts[2]+"\n"
			mp[parts[1]]=parts[2]
			rmu.Unlock()
			mess := "OK\n"
			conn.Write([]byte(mess))
		}else if parts[0]=="get"{
			rmu.RLock()
			value , ok := mp[parts[1]]
			rmu.RUnlock()
			var mess string
			if ok{
				mess="value: "+value+"\n"
				//mess="OK\n"
				//log.Print("req done",mess)
			}else{
				mess="<nil>\n"
			}
			conn.Write([]byte(mess))
		}else if parts[0]=="del"{
			rmu.Lock()
			//log.Print("req is processing [del]")
			var deleted_val=mp[parts[1]]
			delete(mp,parts[1])
			rmu.Unlock()
			mess:="deleted key: "+parts[1]+" with value: "+deleted_val+"\n"
			conn.Write([]byte(mess))
		}else if parts[0]=="inc"{
			rmu.Lock()
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
			rmu.Unlock() 
			mess:="key: "+parts[1]+" increased by "+inc_by+"\n"
			conn.Write([]byte(mess))
		}else if parts[0]=="dec"{
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
		}else{
			conn.Write([]byte("write correct command\n"))
		}
    }
}