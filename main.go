package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"sync"
	"unicode/utf8"
)

var wg sync.WaitGroup

var (
	host     = ""
	password = ""
)

func main() {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go read(conn, &wg)

	conn.Write(makeCmd("AUTH", password))

	conn.Write(makeCmd("SET", "TRY_SET", "THIS IS FOR SET"))
	conn.Write(makeCmd("GET", "TRY_SET"))

	conn.Write(makeCmd("INCR", "TRY_INCR"))
	conn.Write(makeCmd("GET", "TRY_INCR"))

	wg.Wait()
}

func read(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(buffer[:n]))
	}
}

func makeCmd(cmds ...string) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("*")
	buffer.WriteString(strconv.Itoa(len(cmds)))
	for _, cmd := range cmds {
		buffer.WriteString("\r\n")
		buffer.WriteString("$")
		buffer.WriteString(strconv.Itoa(utf8.RuneCountInString(cmd)))
		buffer.WriteString("\r\n")
		buffer.WriteString(cmd)
	}
	buffer.WriteString("\r\n")
	return buffer.Bytes()
}
