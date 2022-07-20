package tests

import (
	"encoding/json"
	model2 "github.com/BlackRRR/nats-streaming/internal/model"
	"io"
	"net/http"
	"testing"
)

type Client struct {
	url string
}

func NewClient(url string) Client {
	return Client{url}
}

func Test_http(t *testing.T) {
	var model model2.Model
	c := NewClient("http://localhost:8011/b563feb7b2b84b6test")
	res, err := c.NewRequest(nil, http.MethodGet)
	if err != nil {
		t.Errorf("expected err to be nil got %v", err)
	}

	err = res.Decode(&model)
	if err != nil {
		t.Errorf("failed to decode response %v", err)
	}

	t.Log(model)
}

func (c Client) NewRequest(r io.Reader, method string) (*json.Decoder, error) {
	client := http.Client{}
	request, err := http.NewRequest(method, c.url, r)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	decode := json.NewDecoder(resp.Body)

	return decode, nil
}
