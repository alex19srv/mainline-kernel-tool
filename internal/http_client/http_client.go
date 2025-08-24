package http_client

import (
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	client http.Client
}

func New() *HTTPClient {
	transport := http.Transport{
		MaxIdleConns:        2,
		MaxIdleConnsPerHost: 2,
		IdleConnTimeout:     5 * time.Second,
	}
	res := &HTTPClient{
		client: http.Client{
			Timeout:   0, // 300 * time.Second,
			Transport: &transport,
		},
	}

	return res
}

func (c *HTTPClient) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
func (c *HTTPClient) GetAndSave(url string, writer io.Writer) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
