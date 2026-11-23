package logg

import (

	"log"
	"bufio"
	"os"
	"path/filepath"
)

type LogFile struct{
	loger 	*log.Logger
	file 	*os.File
	buffer *bufio.Writer
}

const bufferSize =64*1024;

func CreateLogger(path string)(*LogFile,error){
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
    	panic(err)
	}
	
	file, err := os.OpenFile(
        path,
        os.O_APPEND|os.O_CREATE|os.O_WRONLY,
        0666,
    )
    if err != nil {
        panic(err)
    }

	buffer := bufio.NewWriterSize(file, bufferSize)
	loger := log.New(file, "SERVER: ", log.LstdFlags)

	return &LogFile{ file:file , loger:loger , buffer:buffer},nil
}


func (l *LogFile) Close() error {
    err := l.buffer.Flush(); 
	
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

