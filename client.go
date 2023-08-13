package vroom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Url string
}

func NewClient(url string) *Client {
	return &Client{Url: url}
}

func (c *Client) Solve(p *Problem) (*Solution, error) {
	p.processLocations()

	httpClient := http.Client{}
	body, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.Url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body = []byte{}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading body: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("got response %d: %s", resp.StatusCode, string(body))
	}

	result := &Solution{}
	err = json.Unmarshal(body, &result)
	return result, err
}
