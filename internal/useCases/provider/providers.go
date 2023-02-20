package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Providers interface {
	GetResponse() *http.Request
	ParceResponse(*http.Request) int
}

type provider struct {
	url   string
	cur   string
	key   string
	model ModelResponse
}

func NewProvider(url, cur, key string) *provider {
	return &provider{
		url:   url,
		cur:   cur,
		key:   key,
		model: ModelResponse{},
	}
}

func (p *provider) GetResponse() *http.Request {

	URL := fmt.Sprintf(p.url, p.cur)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println(err)
	}

	return req

}

func (p *provider) ParceResponse(req *http.Request) int {

	client := &http.Client{}
	req.Header.Set("apikey", p.key)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ПОЛОМКА В ЗАПРОСЕ")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, _ := io.ReadAll(resp.Body)

	unmarshalErr := json.Unmarshal(body, &p.model)
	if unmarshalErr != nil {
		log.Println(unmarshalErr)
	}

	return int(p.model.Result * 100)
}

func Sender(p Providers) int {

	req := p.GetResponse()
	return p.ParceResponse(req)

}
