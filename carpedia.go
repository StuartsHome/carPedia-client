package carpedia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL = "http://localhost:8100/"
)

type Desc struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

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
	RateLimits [categories]Rate
	RetryDelay time.Duration
	RetryLimit int
	Timeout    time.Duration
}

type Client struct {
	client *http.Client

	mu   sync.Mutex
	opts ClientOpts

	common service

	App  *AppService
	Car  *CarService
	Desc *DescService
}

func NewClient(opts ClientOpts) *Client {
	if opts.Timeout == 0 {
		opts.Timeout = time.Minute
	}
	cl := http.Client{
		Timeout: opts.Timeout,
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	if opts.BaseURL == nil {
		opts.BaseURL = baseURL
	}
	if opts.RetryLimit == 0 {
		opts.RetryLimit = 3
	}
	if opts.RetryDelay == 0 {
		opts.RetryDelay = time.Minute
	}

	c := &Client{
		client: &cl,
		opts:   opts,
	}
	c.common.client = c
	c.App = (*AppService)(&c.common)
	c.Car = (*CarService)(&c.common)
	c.Desc = (*DescService)(&c.common)
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

func (c *Client) defaultDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	/*
		No oauth yet
	*/

	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}

	var resp *http.Response
	var err error

	for i := 0; ; i++ {

		resp, err = c.client.Do(req)
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

	// if resp.Body == nil {
	// 	return nil, errors.New("error empty response body")
	// }

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

type rateLimitCategory uint8

const (
	noRateLimitCheck rateLimitCategory = iota
	rateLimitCheck

	categories
)

// type rates map[string]Rate

// The rate limit for the current client
type Rate struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"` // the number of requests this client can make per hour
	Reset     Timestamp `json: "reset"`
}

type RateLimitError struct {
	Rate     Rate
	Response *http.Response
	Message  string `json:"message"`
	Code     int
}

type Timestamp struct {
	time.Time
}

func (t Timestamp) Equal(u Timestamp) bool {
	return t.Time.Equal(u.Time)
}

func (c *Client) RateLimits(ctx context.Context) (*Rate, *http.Response, error) {
	req, err := c.NewRequest(http.MethodGet, "rate_limit", nil)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println(req)

	ctx = context.WithValue(ctx, noRateLimitCheck, true)
	return nil, nil, nil
}

// TODO: finish up
func (c *Client) checkRateLimitBeforeDo(req *http.Request) *RateLimitError {
	// rate := c.opts.RateLimits
	return nil

}
