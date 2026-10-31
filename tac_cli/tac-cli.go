package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	maxRESPBulkLength  = 64 << 20
	maxRESPArrayLength = 1 << 20
)

type respValue struct {
	kind  byte
	text  string
	items []respValue
	nil   bool
}

func writeCommand(writer *bufio.Writer, parts []string) error {
	if len(parts) == 0 {
		return errors.New("empty command")
	}

	if _, err := fmt.Fprintf(writer, "*%d\r\n", len(parts)); err != nil {
		return err
	}
	for _, part := range parts {
		if _, err := fmt.Fprintf(writer, "$%d\r\n", len([]byte(part))); err != nil {
			return err
		}
		if _, err := writer.WriteString(part); err != nil {
			return err
		}
		if _, err := writer.WriteString("\r\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}

func readRESPLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(line) < 2 || !strings.HasSuffix(line, "\r\n") {
		return "", errors.New("invalid RESP line terminator")
	}
	return line[:len(line)-2], nil
}

func readRESPValue(reader *bufio.Reader, depth int) (respValue, error) {
	if depth > 128 {
		return respValue{}, errors.New("RESP nesting limit exceeded")
	}

	prefix, err := reader.ReadByte()
	if err != nil {
		return respValue{}, err
	}
	switch prefix {
	case '+', '-', ':':
		line, err := readRESPLine(reader)
		if err != nil {
			return respValue{}, err
		}
		return respValue{kind: prefix, text: line}, nil
	case '$':
		line, err := readRESPLine(reader)
		if err != nil {
			return respValue{}, err
		}
		length, err := strconv.Atoi(line)
		if err != nil || length < -1 || length > maxRESPBulkLength {
			return respValue{}, errors.New("invalid RESP bulk string length")
		}
		if length == -1 {
			return respValue{kind: '$', nil: true}, nil
		}

		data := make([]byte, length+2)
		if _, err := io.ReadFull(reader, data); err != nil {
			return respValue{}, err
		}
		if data[length] != '\r' || data[length+1] != '\n' {
			return respValue{}, errors.New("invalid RESP bulk string terminator")
		}
		return respValue{kind: '$', text: string(data[:length])}, nil
	case '*':
		line, err := readRESPLine(reader)
		if err != nil {
			return respValue{}, err
		}
		length, err := strconv.Atoi(line)
		if err != nil || length < -1 || length > maxRESPArrayLength {
			return respValue{}, errors.New("invalid RESP array length")
		}
		if length == -1 {
			return respValue{kind: '*', nil: true}, nil
		}

		value := respValue{kind: '*', items: make([]respValue, length)}
		for i := range value.items {
			value.items[i], err = readRESPValue(reader, depth+1)
			if err != nil {
				return respValue{}, err
			}
		}
		return value, nil
	default:
		return respValue{}, fmt.Errorf("unsupported RESP type %q", prefix)
	}
}

func writeRESPValue(writer io.Writer, value respValue, indent string) {
	switch value.kind {
	case '+':
		fmt.Fprintln(writer, value.text)
	case '$':
		if value.nil {
			fmt.Fprintln(writer, "(nil)")
			return
		}
		fmt.Fprintln(writer, value.text)
	case '-':
		fmt.Fprintf(writer, "(error) %s\n", value.text)
	case ':':
		fmt.Fprintln(writer, value.text)
	case '*':
		if value.nil {
			fmt.Fprintln(writer, "(nil)")
			return
		}
		for i, item := range value.items {
			fmt.Fprintf(writer, "%s%d) ", indent, i+1)
			if item.kind == '*' {
				fmt.Fprintln(writer)
				writeRESPValue(writer, item, indent+"  ")
			} else if item.kind == '-' {
				fmt.Fprintf(writer, "(error) %s\n", item.text)
			} else if item.kind == '$' && item.nil {
				fmt.Fprintln(writer, "(nil)")
			} else {
				fmt.Fprintln(writer, item.text)
			}
		}
	}
}

func parseCommand(line string) ([]string, error) {
	var parts []string
	var current strings.Builder
	var quote rune
	escaped := false
	inToken := false

	for _, char := range line {
		if escaped {
			current.WriteRune(char)
			escaped = false
			inToken = true
			continue
		}
		if char == '\\' {
			escaped = true
			inToken = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				current.WriteRune(char)
			}
			continue
		}
		switch {
		case char == '"' || char == '\'':
			quote = char
			inToken = true
		case char == ' ' || char == '\t':
			if inToken {
				parts = append(parts, current.String())
				current.Reset()
				inToken = false
			}
		default:
			current.WriteRune(char)
			inToken = true
		}
	}

	if escaped {
		return nil, errors.New("command ends with an incomplete escape")
	}
	if quote != 0 {
		return nil, errors.New("unterminated quoted argument")
	}
	if inToken {
		parts = append(parts, current.String())
	}
	return parts, nil
}

func execute(reader *bufio.Reader, writer *bufio.Writer, parts []string) error {
	if err := writeCommand(writer, parts); err != nil {
		return fmt.Errorf("send command: %w", err)
	}
	value, err := readRESPValue(reader, 0)
	if err != nil {
		return fmt.Errorf("read server response: %w", err)
	}
	writeRESPValue(os.Stdout, value, "")
	return nil
}

func run() error {
	flags := flag.NewFlagSet("tac-cli", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	host := flags.String("h", "127.0.0.1", "server host")
	flags.StringVar(host, "host", "127.0.0.1", "server host")
	port := flags.Int("p", 6379, "server port")
	flags.IntVar(port, "port", 6379, "server port")
	if err := flags.Parse(os.Args[1:]); err != nil {
		return err
	}
	if *port < 1 || *port > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", *port)
	}

	address := net.JoinHostPort(*host, strconv.Itoa(*port))
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", address, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	args := flags.Args()
	if len(args) > 0 {
		return execute(reader, writer, args)
	}

	info, err := os.Stdin.Stat()
	if err != nil {
		return fmt.Errorf("inspect stdin: %w", err)
	}
	interactive := info.Mode()&os.ModeCharDevice != 0
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	for {
		if interactive {
			fmt.Fprintf(os.Stdout, "%s> ", address)
		}
		if !scanner.Scan() {
			break
		}
		parts, err := parseCommand(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERR %v\n", err)
			continue
		}
		if len(parts) == 0 {
			continue
		}
		switch strings.ToLower(parts[0]) {
		case "quit", "exit":
			return nil
		}
		if err := execute(reader, writer, parts); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read command: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
