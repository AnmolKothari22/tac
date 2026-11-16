package client

import (
	"strings"
	"fmt"
)

func encodeCommand(args... string)string{
	var s strings.Builder
	 fmt.Fprintf(&s, "*%d\r\n", len(args))
    for _, arg := range args {
        fmt.Fprintf(&s, "$%d\r\n%s\r\n", len(arg), arg)
    }
    return s.String()
}