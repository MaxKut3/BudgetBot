package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {

	var token string = "5859735143:AAGFvrZyYDwF7RlOkrca4Ekfx8rPpJJQU9k"
	var url = "https://api.telegram.org/bot" + token

	for {
		ans, err := getUpdate(url)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(ans)
	}

}

func getUpdate(url string) ([]UpdateModel, error) {

	update, err := http.Get(url + "/getUpdates")
	if err != nil {
		return nil, err
	}

	defer update.Body.Close()

	date, err := io.ReadAll(update.Body)
	if err != nil {
		return nil, err
	}

	var res Request

	if err := json.Unmarshal(date, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}
