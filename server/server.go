package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleIncomingRequest(conn)
	}
}
func handleIncomingRequest(conn net.Conn) {

	var resultData string

	// Принимаем данные от клиента
	message, _ := bufio.NewReader(conn).ReadString(' ')
	// Убираем лишние пробелы
	message1 := strings.ReplaceAll(message, "\r\n", ",")
	// Разделяем числа по запятой
	message2 := strings.Split(message1, ",")

	// Идем по массиву строки
	for i := 0; i < len(message2)-2; i += 2 {
		if message2[i] != " " {
			// Вытаскиваем числа и перемножаем
			num1, _ := strconv.Atoi(message2[i])
			num2, _ := strconv.Atoi(message2[i+1])
			num3 := strconv.Itoa(num1 * num2)
			// Записываем все в итоговую строку для отправки
			resultData += num3 + "," + "\r\n"
		}
	}

	resultData += "\r\n"

	conn.Write([]byte(resultData))

	conn.Close()
}
