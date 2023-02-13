package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func fixerAPI(cur, key string) int {

	URL := fmt.Sprintf("https://api.apilayer.com/fixer/convert?to=RUB&from=%s&amount=1", cur)

	client := &http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("apikey", key)

	res, err := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)

	var fixerJSON FixerJSON

	unmarshalErr := json.Unmarshal(body, &fixerJSON)
	if unmarshalErr != nil {
		log.Println("fixerAPI не распарсился")
		log.Println(unmarshalErr)
	}

	return int(fixerJSON.Result * 100)
}

func coinGateAPI(cur, key string) int {

	url := fmt.Sprintf("https://api.coingate.com/api/v2/rates/merchant/%s/RUB", cur)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "text/plain")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	exchangeRate, _ := strconv.ParseFloat(string(body), 64)
	return int(exchangeRate * 100)
}

func exchangeratesAPI(cur, key string) int {

	url := fmt.Sprintf("https://api.apilayer.com/exchangerates_data/convert?to=RUB&from=%s&amount=1", cur)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", key)

	if err != nil {
		fmt.Println(err)
	}
	res, _ := client.Do(req)
	defer res.Body.Close()

	var exchangeratesJSON ExchangeratesJSON

	body, _ := io.ReadAll(res.Body)

	errUnmarsh := json.Unmarshal(body, &exchangeratesJSON)
	if errUnmarsh != nil {
		log.Println("exchangeratesAPI не распарсился")
		log.Println(errUnmarsh)
	}

	return int(exchangeratesJSON.Result * 100)
}
