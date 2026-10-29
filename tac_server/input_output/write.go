package inputoutput
import(
	"net"
	"fmt"
)

func Write_RESP_bulk(s string)string{
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func Write_RESP_ERR(s string)string{
	 return fmt.Sprintf("-ERR %s\r\n", s)
}

func MissingKey() string {
	return "$-1\r\n"
}


func Write_String(conn net.Conn,mess string){
	mess = mess+"\n"
	conn.Write([]byte(mess))
}