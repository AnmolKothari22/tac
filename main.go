package main

import (
	
	"fmt"
//	"log"
	"net"
	"os"
	"os/signal"

	"sync"
	"syscall"
	"tac/handle"
//	"time"
)




func main(){
	var wg sync.WaitGroup
    var mu sync.Mutex
	var rmu sync.RWMutex

	lstr , err := net.Listen("tcp",":6379")
	if(err != nil){
		fmt.Print("error occured: ",err,"\n")
		return
	}

	sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
        <-sig
        lstr.Close()
    }()

	mp := make(map[string]string)
	for {
		conn ,err2:= lstr.Accept()
		if(err2!=nil){
			fmt.Print("error : ",err2)
			break
		}
		wg.Add(1)
		go handle.Handle(conn,mp,&wg,&mu,&rmu)
	}
	wg.Wait()
	fmt.Println("server close");


}