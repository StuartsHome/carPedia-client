package carpedia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "http://localhost:8100/"
)

type Response struct {
	*http.Response
}

type service struct {
	client *Client
}

type ClientOpts struct {
	ApiKey     string
	ApiSecret  string
	BaseURL    *url.URL
	RetryLimit int
	RetryDelay time.Duration
	Timeout    time.Duration
}

type Client struct {
	client *http.Client
	opts   ClientOpts

	common service

	Car *CarService
}

func NewClient(httpClient *http.Client, opts ClientOpts) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client: httpClient,
		opts: ClientOpts{
			BaseURL: baseURL,
		},
	}
	c.common.client = c
	c.Car = (*CarService)(&c.common)

	return c

}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// return c.defaultDo(ctx, req)
	return c.client.Do(req)

}

// func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {

// }

type APIError struct {
	Code    int    `json: "code"`
	Message string `json: "message"`
}

func (c *Client) defaultDo(ctx context.Context, req *http.Request) (*Response, error) {
	/*
		No oauth yet
	*/

	var resp *Response
	var err error
	for i := 0; ; i++ {

		// resp, err := c.Do(ctx, req)
		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusTooManyRequests {
			break
		}
		if i >= c.opts.RetryLimit {
			break
		}
		time.Sleep(c.opts.RetryDelay)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiErr APIError
	err = json.Unmarshal(body, &apiErr)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}
	return resp, nil

}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.opts.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have trailing slash, but %q does not", c.opts.BaseURL)
	}

	// parse BaseURL
	u, err := c.opts.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// encode body into io.ReadWriter
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		json.NewEncoder(buf).Encode(&body)
	}

	// create request
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// TODO: userAgent?
	return req, nil

}
