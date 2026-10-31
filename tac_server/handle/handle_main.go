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
	"github.com/AnmolKothari22/tac/config"
	"github.com/AnmolKothari22/tac/tac_server/commands"
	"github.com/AnmolKothari22/tac/tac_server/hasher"
	"github.com/AnmolKothari22/tac/tac_server/input_output"
	"github.com/AnmolKothari22/tac/tac_server/shards"
	//	"time"
)

func handle_logic(conn net.Conn, mp []shards.Shard, parts []string) {
	if len(parts) == 0 {
		conn.Write([]byte(inputoutput.Write_RESP_ERR("empty command")))
		return
	}

	parts[0] = strings.ToLower(parts[0])
	switch parts[0] {
	case "set":
		if len(parts) != 3 {
			conn.Write([]byte(inputoutput.Write_RESP_ERR("wrong number of arguments for 'set' command")))
			return
		}
	case "get", "del":
		if len(parts) != 2 {
			conn.Write([]byte(inputoutput.Write_RESP_ERR("wrong number of arguments for '" + parts[0] + "' command")))
			return
		}
	case "inc", "dec":
		if len(parts) < 2 || len(parts) > 3 {
			conn.Write([]byte(inputoutput.Write_RESP_ERR("wrong number of arguments for '" + parts[0] + "' command")))
			return
		}
	default:
		conn.Write([]byte(inputoutput.Write_RESP_ERR("unknown command '" + parts[0] + "'")))
		return
	}

	shardNo := hasher.ShardFor(parts[1], config.Shard_count)
	mapTemp := &mp[shardNo]
	switch parts[0] {
	case "set":
		commands.Handle_set(&mapTemp.Rmu, parts, conn, &mapTemp.Store)
	case "get":
		commands.Handle_get(&mapTemp.Rmu, parts, conn, &mapTemp.Store)
	case "del":
		commands.Handle_del(&mapTemp.Rmu, parts, conn, &mapTemp.Store)
	case "inc":
		commands.Handle_inc(&mapTemp.Rmu, parts, conn, &mapTemp.Store)
	case "dec":
		commands.Handle_dec(&mapTemp.Rmu, parts, conn, &mapTemp.Store)
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
