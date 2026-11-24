package logg

import (

	"log"
	"bufio"
	"os"
	"path/filepath"
)

type LogFile struct{
	loger 	*log.Logger //for writing into file
	file 	*os.File   //file object(use to interact with it)
	buffer *bufio.Writer // buffer for writing into file (makes it faster)
}

const bufferSize =64*1024; // x*1024--> KB , x*1024*1024--> MB 

func CreateLogger(path string)(*LogFile,error){
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil { // making directory log if it doesnt exist
    	panic(err)
	}
	
	file, err := os.OpenFile( // opening file
        path,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY,
        0666,
    )
    if err != nil {
        panic(err)
    }

	buffer := bufio.NewWriterSize(file, bufferSize)  
	loger := log.New(file, "SERVER: ", log.LstdFlags)

	return &LogFile{ file:file , loger:loger , buffer:buffer},nil // returning pointer to object created
}


func (l *LogFile) Close() error {
    err := l.buffer.Flush(); // flushing buffer before closing 
	
	if err != nil {
        l.file.Close()
        return err
    }
    return l.file.Close()
}

func (l *LogFile)WriteString(s ... string ){
	l.loger.Println(s);
}

func (l *LogFile)WriteStringArray(arr[]string){
	l.loger.Println(arr);
}

