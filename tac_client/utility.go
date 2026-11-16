
package client

import (
    "bufio"
    "fmt"
    "io"
    "strconv"
    "strings"
)

func (c *Client) execute(args ...string) (string, error) {
    c.mu.Lock()
    defer c.mu.Unlock()
    _, err := c.conn.Write([]byte(encodeCommand(args...))) //writing in resp to server
    if err != nil {
        return "", err
    }

    return readBulkString(c.reader)
}


func readBulkString(r *bufio.Reader) (string, error) {
    header, err := r.ReadString('\n') // seperates header(lenght of whole argument array part ) eg : header-> "$5\r\n"
    if err != nil {
        return "", err
    }
	
    header = strings.TrimSuffix(header, "\r\n") //trimming suffix from header part

    if strings.HasPrefix(header, "-") { //error handling
        return "", fmt.Errorf("%s", header[1:])
    }

    if header == "$-1" { //missing key
        return "", fmt.Errorf("key not found")
    }

    if !strings.HasPrefix(header, "$") {  //if string DOESNT start with $
        return "", fmt.Errorf("unexpected RESP response: %s", header)
    }

    n, err := strconv.Atoi(header[1:])  //extracting lenght from header 
    if n < 0 {
        return "", fmt.Errorf("invalid bulk string length")
    }
    if err != nil {
        return "", err
    }

    data := make([]byte, n+2) // making space for n lenght string and \r\n

    if _, err := io.ReadFull(r, data); err != nil {
        return "", err
    }

    if string(data[n:]) != "\r\n" {
        return "", fmt.Errorf("invalid bulk string terminator")
    }

    return string(data[:n]), nil
}
