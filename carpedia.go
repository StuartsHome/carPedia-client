package carpedia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	RateLimits rates
	RetryLimit int
	RetryDelay time.Duration
	Timeout    time.Duration
}

type Client struct {
	client *http.Client
	opts   ClientOpts

	common service

	App  *AppService
	Car  *CarService
	Desc *DescService

	// do we need a mutex?
}

func NewClient(opts ClientOpts) *Client {
	cl := http.Client{Timeout: time.Minute}
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client: &cl,
		opts: ClientOpts{
			BaseURL: baseURL,
		},
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

func (c *Client) defaultDo(ctx context.Context, req *http.Request) (*Response, error) {
	/*
		No oauth yet
	*/

	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}

	// if the current request is subject to rate limiting
	if rateLimits := ctx.Value(noRateLimitCheck); rateLimits == nil {
		// if we've hit the rate limit
		// if err := c
	}

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

const noRateLimitCheck = iota

type rates map[string]Rate

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
