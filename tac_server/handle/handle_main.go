package handle

import (
	"bufio"
	//"fmt"
	//"log"
	"net"
	"strings"
	"sync"

	//"strconv"
	//"io"
	"github.com/AnmolKothari22/tac"
	"github.com/AnmolKothari22/tac/tac_server/commands"
	"github.com/AnmolKothari22/tac/tac_server/hasher"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
	//	"time"
)



func handle_logic(conn net.Conn, mp []shards.Shard, parts []string){
	//lowercasing first input
		//fmt.Print("stage 4");
		parts[0]=strings.ToLower(parts[0])
		//getting shard number
		shard_no := hasher.ShardFor(parts[1],tac.Shard_count);
		//fmt.Print("--",shard_no,"\n")
		map_temp := &mp[shard_no];
		if parts[0]=="set"{
			commands.Handle_set(&map_temp.Rmu,parts,conn,&map_temp.Store);
		}else if parts[0]=="get"{
			commands.Handle_get(&map_temp.Rmu,parts,conn,&map_temp.Store);
		}else if parts[0]=="del"{
			commands.Handle_del(&map_temp.Rmu,parts,conn,&map_temp.Store);
		}else if parts[0]=="inc"{
			commands.Handle_inc(&map_temp.Rmu,parts,conn,&map_temp.Store);
		}else if parts[0]=="dec"{
			commands.Handle_dec(&map_temp.Rmu,parts,conn,&map_temp.Store);
		}else{
			conn.Write([]byte("write correct command\n"))
		}

}

func Handle(conn net.Conn, mp []shards.Shard,wg *sync.WaitGroup,mode int) {
	defer conn.Close()
	defer wg.Done()
	//fmt.Print("stage 1");
	//RESP
	if mode==1{
		reader := bufio.NewReader(conn)
		//fmt.Print("stage 2");
		for {
    		parts, err := inputoutput.ReadRESP(reader)
    		if err != nil {
				//log.Println("error 2: ",err);
        		break
    		}	
			//fmt.Print("stage 3");
			handle_logic(conn,mp,parts);
    	}
	}
	//non RESP
	if mode==2{
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
        	parts := strings.Fields(scanner.Text())
			handle_logic(conn,mp,parts);
		}
	}
}