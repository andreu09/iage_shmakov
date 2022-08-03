package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

// Структура данных запроса Test1
type DataTest1 struct {
	Key string `json:"key"`
	Val int64  `json:"val"`
}

// Структура данных запроса Test2
type DataTest2 struct {
	S   string `json:"s"`
	Key string `json:"key"`
}

// Структура данных запроса  Test3
type DataTest3 []struct {
	A   string `json:"a"`
	B   string `json:"b"`
	Key string `json:"key"`
}

// Структура данных результата Test3
type DataTest3Result struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Подключение к Redis
func ConnectRedis(host string, port string) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     host + port,
		Password: "",
		DB:       0,
	})

}

func main() {

	host := flag.String("host", "localhost:", "Хост")
	port := flag.String("port", "6379", "Порт")
	flag.Parse()

	ConnectRedis(*host, *port)

	http.HandleFunc("/test1", Test1)
	http.HandleFunc("/test2", Test2)
	go http.HandleFunc("/test3", Test3)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}

/*
curl --header "Content-Type: application/json" --request POST --data '{"key":"t2est","val":12}' http://localhost:8000/test1
Обработка запросов по пути /test1
*/
func Test1(w http.ResponseWriter, r *http.Request) {

	// Данные с запроса
	var response DataTest1

	// Чтение запроса
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		// Ошибка запроса
		fmt.Println(w, "err %q\n", err, err.Error())
	} else {
		// Если ошибок нет, записываем даные по указателю в структуру
		err = json.Unmarshal(body, &response)
		if err != nil {
			// Не удалось раскодировать данные json
			fmt.Println(w, "can't unmarshal: ", err.Error())
		}
	}

	if err == redis.Nil {
		// Ключ раннее не существовал
		rdb.Set(response.Key, response.Val, 0).Err()
	} else if err != nil {
		// Произошла ошибка при записи ключа
		panic(err)
	} else {
		// Инкремент ключа
		pipe := rdb.TxPipeline()
		incr := pipe.IncrBy(response.Key, response.Val)
		_, exer := pipe.Exec()
		if exer != nil {
			// Ошибка при инкрементации
			fmt.Println(exer)
		}

		resultData := map[string]string{response.Key: strconv.Itoa(int(incr.Val()))}
		resultJsonData, _ := json.Marshal(resultData)
		io.WriteString(w, string(resultJsonData))

	}
}

/*
	curl --header "Content-Type: application/json" --request POST --data '{"key":"t2est","val":12}' http://localhost:8000/test2
	Обработка запросов по пути /test2
*/

func Test2(w http.ResponseWriter, r *http.Request) {

	// Данные с запроса
	var response DataTest2
	// Чтение запроса
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		// Ошибка запроса
		fmt.Println(w, "err %q\n", err, err.Error())
	} else {
		// Если ошибок нет, записываем даные по указателю в структуру
		err = json.Unmarshal(body, &response)
		if err != nil {
			// Не удалось раскодировать данные json
			fmt.Println(w, "can't unmarshal: ", err.Error())
		}
	}

	sha_512 := sha512.New()
	sha_512.Write([]byte(response.S))
	hmac512 := hmac.New(sha512.New, []byte(response.Key))
	hmac512.Write([]byte(response.S))

	resultData := base64.StdEncoding.EncodeToString(hmac512.Sum(nil))

	io.WriteString(w, resultData)
}

/*
	curl --header "Content-Type: application/json" --request POST --data '[ {"a":"12", "b":"43", "key":"x"}, {"a":"11", "b":"33", "key":"y"}]' http://localhost:8000/test3
	Обработка запросов по пути /test3
*/

func Test3(w http.ResponseWriter, r *http.Request) {

	// Данные с запроса
	var response DataTest3
	// Строка отправляемая на сервер
	var str string

	// Чтение запроса
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		// Ошибка запроса
		fmt.Println(w, "err %q\n", err, err.Error())
	} else {
		// Если ошибок нет, записываем даные по указателю в структуру
		err = json.Unmarshal(body, &response)
		if err != nil {
			// Не удалось раскодировать данные json
			fmt.Println(w, "can't unmarshal: ", err.Error())
		}
	}

	// Заполняем строку согласно виду из задания
	for _, value := range response {
		str += value.A + "," + value.B + "\r\n"
	}

	str += "\r\n"

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	// Отправляем строку на сервер
	fmt.Fprintf(conn, str+" ")
	// Слушаем ответ от сервера
	message, _ := bufio.NewReader(conn).ReadString(' ')
	// Удаляем лишние пробелы
	messageTrim := strings.TrimSpace(message)
	// Удаляем лишние спецсимволы
	messageTrimReplace := strings.Replace(messageTrim, "\r\n", "", -1)
	// Разделяем числа по запятой
	messageResult := strings.Split(messageTrimReplace, ",")

	x, _ := strconv.Atoi(messageResult[0])
	y, _ := strconv.Atoi(messageResult[1])

	resultResponse := DataTest3Result{x, y}
	resultJsonData, _ := json.Marshal(resultResponse)

	io.WriteString(w, string(resultJsonData))

}
