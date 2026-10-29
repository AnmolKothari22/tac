package main

import (
	"flag"
	"fmt"

	//	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	
	"github.com/AnmolKothari22/tac/tac_server/handle"
	"github.com/AnmolKothari22/tac/tac_server/shards"
	//	"time"
)




const mode=1;

func main(){
	var wg sync.WaitGroup
	//var rmu sync.RWMutex
	//for running on any ports 
	portPtr :=flag.String("port","6379","port where tac server has to run")
	flag.Parse()
	PORT :=":"+(*portPtr)
	//starting server 
	lstr , err := net.Listen("tcp",PORT)
	if(err != nil){
		fmt.Print("error occured: ",err,"\n")
		return
	}

	//graceful shutdown 
	sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
        <-sig
        lstr.Close()
    }()
	//main map
//	mp := make(map[string]string)
	map_storage := make([]shards.Shard, shards.Shard_num())
	//initializing maps inside each struct
	for i := range map_storage {
    	map_storage[i].Store = make(map[string]string)
	}
	for {
		conn ,err2:= lstr.Accept()
		if(err2!=nil){
			fmt.Print("error 1 : ",err2)
			break
		}
		wg.Add(1)
		go handle.Handle(conn,map_storage,&wg,mode)
	}
	wg.Wait()
	fmt.Println("server close");


}

/*

mkdir myapp && cd myapp
go mod init github.com/AnmolKothari22/tac



*/