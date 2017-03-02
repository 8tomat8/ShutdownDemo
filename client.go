package main

import (
	"bufio"
	"net/http"
	"fmt"
)

func main() {
	resp, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		panic(err)
	}

	count := 1
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}

		fmt.Println(count, string(line))
		count++
	}
}
