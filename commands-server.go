package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

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
		// // Будем прослушивать все сообщения разделенные \n
		// message, err := bufio.NewReader(conn).ReadString('\n')
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// // Распечатываем полученое сообщение
		// fmt.Print("Message Received:", string(message))
		// // Процесс выборки для полученной строки
		// // newmessage := strings.ToUpper(message)
		// // Отправить новую строку обратно клиенту
		// conn.Write([]byte(time.Now().Format(time.RFC3339) + "\n"))
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
	message := string(bufferBytes)
	// get the remote address of the client
	clientAddr := conn.RemoteAddr().String()
	// format a response
	response := fmt.Sprintf(message + " from " + clientAddr + "\n")

	// have server print out important information
	log.Println(response)

	// let the client know what happened
	// conn.Write([]byte("you sent: " + response))
	conn.Write([]byte(time.Now().Format(time.RFC3339) + response + "\n"))

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn)
}
