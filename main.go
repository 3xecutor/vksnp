package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ApiResponse struct {
	Response json.RawMessage `json:"response"`
}

func main() {
	inputFile, err := os.Open("shortname.txt")
	if err != nil {
		log.Fatalf("Не удалось открыть файл для чтения: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("Не удалось открыть файл для записи: %v", err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()

		data := url.Values{}
		data.Set("access_token", "vk1.a.-DF3Q") // замените на реальный токен
		data.Set("screen_name", line)           // замените на реальный screen name
		data.Set("v", "5.888")                  // замените на реальный screen name

		// Создаем запрос к API
		resp, err := http.PostForm("https://api.vk.com/method/utils.resolveScreenName", data)
		if err != nil {
			log.Printf("Ошибка при вызове API для строки %s: %v", line, err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Ошибка при чтении ответа от API для строки %s: %v", line, err)
			continue
		}

		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			log.Printf("Ошибка при разборе JSON ответа для строки %s: %v", line, err)
			continue
		}

		var responseArray []interface{}
		if err = json.Unmarshal(apiResponse.Response, &responseArray); err == nil {
			if len(responseArray) == 0 {
				_, err := outputFile.WriteString(line + "\n")
				if err != nil {
					log.Printf("Ошибка при записи строки в файл: %v", err)
				}
			}
		}
	}

	time.Sleep(1 * time.Second)

	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}
}
