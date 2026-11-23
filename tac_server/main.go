package main

import (
	"flag"
	"fmt"
	//"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/AnmolKothari22/tac/config"
	"github.com/AnmolKothari22/tac/tac_server/handle"
	"github.com/AnmolKothari22/tac/tac_server/shards"
	"github.com/AnmolKothari22/tac/tac_server/log_handler"
	//"time"
	//"log"
)




const mode=config.Mode;

func main(){
	var wg sync.WaitGroup

	//for running on any ports 
	portPtr :=flag.String("port",strconv.Itoa(config.PORT),"port where tac server has to run")
	flag.Parse()
	PORT :=":"+(*portPtr)
	//starting server 
	lstr , err := net.Listen("tcp",PORT)
	if(err != nil){
		fmt.Print("error occured: ",err,"\n")
		return
	}

	//logger created
	logger ,err :=logg.CreateLogger("./tac_server/log_handler/logs/ServerInfo.log")
	defer logger.Close()
	logger.WriteString("STARTING SERVER...")
	logger.WriteString(`Server started at `,PORT)
	
	//graceful shutdown 
	sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
        <-sig
        lstr.Close()
    }()
	
	//main map
	map_storage := make([]shards.Shard, config.Shard_count)
	//initializing maps inside each struct
	for i := range map_storage {
    	map_storage[i].Store = make(map[string]shards.Stored_values)
	}
	for {
		conn ,err2:= lstr.Accept()
		if(err2!=nil){
			fmt.Print("error 1 : ",err2)
			break
		}
		wg.Add(1)
		go handle.Handle(conn,map_storage,&wg,logger,mode)
	}
	wg.Wait()
	fmt.Println("server close");


}

/*

loger_server := logg.Create_Log("")

okk just test
test2
*/