package main

import (
	
	"fmt"
	"flag"
//	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tac_server/handle"
//	"time"
)




func main(){
	var wg sync.WaitGroup
	var rmu sync.RWMutex
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
	mp := make(map[string]string)
	for {
		conn ,err2:= lstr.Accept()
		if(err2!=nil){
			fmt.Print("error : ",err2)
			break
		}
		wg.Add(1)
		go handle.Handle(conn,mp,&wg,&rmu)
	}
	wg.Wait()
	fmt.Println("server close");


}