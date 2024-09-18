package main

import (
	"fmt"
	"io"
	"net/http"
)

func downloadURL(url string) ([]byte, error) {
	// Получаем ответ от get запроса на указанный url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Если ответ не 200, выдаём ошибку
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP response error: %s", resp.Status)
	}

	// Читаем ответ в переменную data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
