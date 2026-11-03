package inputoutput
import(
	"net"
	"fmt"
)
type WriteType int
const(
	Normal WriteType = iota
	Error
	KeyMiss
)

func Write_RESP_bulk(s string)string{ //1
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func Write_RESP_ERR(s string)string{ //2
	 return fmt.Sprintf("-ERR %s\r\n", s)
}

func MissingKey() string { //3
	return "$-1\r\n"
}


func Write_String(conn net.Conn,mess string){
	mess = mess+"\n"
	conn.Write([]byte(mess))
}