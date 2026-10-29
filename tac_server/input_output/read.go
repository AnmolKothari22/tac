package inputoutput

import (
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"io"
)

func ReadRESP(reader *bufio.Reader) ([]string, error) {
	// Read array header: *3\r\n
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
        	return nil, err
    	}
		return nil, err
	}

	if !strings.HasSuffix(line, "\r\n") || len(line) < 3 || line[0] != '*' {
		return nil, fmt.Errorf("invalid RESP array header")
	}

	count, err := strconv.Atoi(line[1 : len(line)-2])
	if err != nil || count < 1 || count > 1024 {
		return nil, fmt.Errorf("invalid RESP array length")
	}

	parts := make([]string, 0, count)

	for i := 0; i < count; i++ {
		// Read bulk string header: $3\r\n
		line, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("2",err);
			return nil, err
		}

		if !strings.HasSuffix(line, "\r\n") ||
			len(line) < 3 || line[0] != '$' {
			return nil, fmt.Errorf("invalid RESP bulk string header")
		}

		length, err := strconv.Atoi(line[1 : len(line)-2])
		if err != nil || length < 0 || length > 1<<20 {
			return nil, fmt.Errorf("invalid RESP bulk string length")
		}

		// Read exactly 'length' bytes, then the trailing CRLF.
		data := make([]byte, length+2)
		if _, err := io.ReadFull(reader, data); err != nil {
			fmt.Println("3",err);
			return nil, err
		}

		if data[length] != '\r' || data[length+1] != '\n' {
			return nil, fmt.Errorf("missing bulk string terminator")
		}

		parts = append(parts, string(data[:length]))
	}

	return parts, nil
}

/*



*/