package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

type Command interface {
	Run() string
}

type GetTime struct {
}

func (command GetTime) Run() string {
	return time.Now().Format(time.RFC3339)
}

type GetCPUusage struct {
}

func (command GetCPUusage) Run() string {
	percent, _ := cpu.Percent(time.Second, true)
	return fmt.Sprintf("User CPU usage %.2f", percent[cpu.CPUser])
}

type DefaultCommand struct {
}

func (command DefaultCommand) Run() string {
	return "Type get_cpu or get_time. I'm not smart!"
}

func CommandDispatcher(command_name string) Command {
	switch command_name {
	case "get_time":
		return GetTime{}
	case "get_cpu":
		return GetCPUusage{}
	default:
		return DefaultCommand{}
	}
}

// требуется только ниже для обработки примера

func main() {
	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, err := net.Listen("tcp", ":8881")
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}

	for {
		// Открываем порт
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Будем прослушивать все сообщения разделенные \n
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Println("client left..")
		conn.Close()

		// escape recursion
		return
	}

	// convert bytes from buffer to string
	message := strings.TrimSuffix(string(bufferBytes), "\n")
	// get the remote address of the client
	clientAddr := conn.RemoteAddr().String()
	// format a response
	response := fmt.Sprintf(message + " from " + clientAddr + "\n")

	// have server print out important information
	log.Println(response)

	// conn.Write([]byte(GetTime{}.Run() + "\n"))
	conn.Write([]byte(CommandDispatcher(message).Run() + "\n"))

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn)
}
