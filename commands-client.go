package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func exitCatcher(conn net.Conn) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(conn net.Conn) {
		<-sigs
		fmt.Println("\nExit...")
		conn.Close()
		os.Exit(1)
	}(conn)
}

func main() {
	// Подключаемся к сокету
	conn, err := net.Dial("tcp", "127.0.0.1:8881")
	if err != nil {
		log.Fatal(err)
	}

	go exitCatcher(conn)

	for {
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		fmt.Print(text)
		// Отправляем в socket
		fmt.Fprintf(conn, text+"\n")
		// Прослушиваем ответ
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal("Server response error: ", err)
		}
		fmt.Print("Message from server: " + message)
	}
}
