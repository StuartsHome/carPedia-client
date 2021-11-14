package carpedia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultBaseURL      = "http://localhost:8100/"
	headerRateLimit     = "X-RateLimit-Limit"     // The maximum number of requests you're permitted to make per hour.
	headerRateRemaining = "X-RateLimit-Remaining" // The number of requests remaining in the current rate limit window.
	headerRateReset     = "X-RateLimit-Reset"     // The time at which the current rate limit window resets in UTC seconds
)

type Desc struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Response struct {
	*http.Response

	Rate Rate
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.Rate = parseRate(r)
	// TODO: timeout token here
	return response
}

func parseRate(r *http.Response) Rate {
	var rate Rate
	if limit := r.Header.Get(headerRateLimit); limit != "" {
		rate.Limit, _ = strconv.Atoi(limit)
	}
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		if v, _ := strconv.ParseInt(reset, 10, 64); v != 0 {
			rate.Reset = Timestamp{time.Unix(v, 0)}
		}
	}
	return rate
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

func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.defaultDo(ctx, req)
	// return c.client.Do(req)

}

type APIError struct {
	Code     int    `json: "code"`
	Message  string `json: "message"`
	Response *http.Response
}

func (a *APIError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		a.Response.Request.Method, a.Response.Request.URL,
		a.Response.StatusCode, a.Message)

}

func (c *Client) defaultDo(ctx context.Context, req *http.Request) (*Response, error) {
	/*
		No oauth yet
	*/

	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}

	var resp *Response
	var err error

	for i := 0; ; i++ {

		resp, err = c.Do(ctx, req)
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

	return resp, nil

}

func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("empty context")
	}

	// get current category
	category := rateLimitCheck

	// check rate limits
	if err := c.checkRateLimitBeforeDo(req, category); err != nil {
		return &Response{
			Response: err.Response,
			Rate:     err.Rate,
		}, nil
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	r := newResponse(response)

	// get rates
	c.mu.Lock()
	c.opts.RateLimits[category] = r.Rate
	c.mu.Unlock()

	if err := c.checkResponse(response); err != nil {
		defer response.Body.Close()
		return nil, err
	}

	return r, nil
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

// Error method to allow RateLimitError to implement the Error interface
func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, formateRateReset(time.Until(r.Rate.Reset.Time)))
}

func formateRateReset(d time.Duration) string {
	// isNegative := d < 0
	// if isNegative {
	// 	d *= -1
	// }
	// secondsTotal := int(0.5 + d.Seconds())
	secondsTotal := int(0.5 + math.Abs(d.Seconds()))
	minutes := secondsTotal / 60
	seconds := secondsTotal - minutes*60

	var timeString string
	if minutes > 0 {
		timeString = fmt.Sprintf("%dm%02ds", minutes, seconds)
	} else {
		timeString = fmt.Sprintf("%ds", seconds)
	}

	if d < 0 {
		return fmt.Sprintf("[rate limit was reset %v ago]", timeString)
	}
	return fmt.Sprintf("[rate reset in %v]", timeString)

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
func (c *Client) checkRateLimitBeforeDo(req *http.Request, category rateLimitCategory) *RateLimitError {
	c.mu.Lock()
	rate := c.opts.RateLimits[category]
	c.mu.Unlock()

	if rate.Reset.Time.IsZero() && rate.Remaining == 0 && time.Now().Before(rate.Reset.Time) {
		resp := &http.Response{
			Status:     http.StatusText(http.StatusForbidden),
			StatusCode: http.StatusForbidden,
			Request:    req,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		return &RateLimitError{
			Rate:     rate,
			Response: resp,
			Message:  fmt.Sprintf("API rate limit of %v still exceeded until %v, not making request.", rate.Limit, rate.Reset.Time),
		}
	}
	return nil
}

// TODO: finish up
func (c *Client) checkResponse(r *http.Response) error {
	// if request contains statuscode in 200 range
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &APIError{Response: r}
	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	switch {
	case r.StatusCode == http.StatusForbidden && r.Header.Get(headerRateRemaining) == "0":
		return &RateLimitError{
			Rate:     parseRate(r),
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
			Code:     r.StatusCode,
		}
	case r.StatusCode == http.StatusForbidden:
		if v := r.Header["Retry-After"]; len(v) > 0 {
			retryAfterSeconds, _ := strconv.ParseInt(v[0], 10, 64)
			retryAfter := time.Duration(retryAfterSeconds) * time.Second
			return &RateLimitError{
				Rate:     parseRate(r),
				Response: errorResponse.Response,
				Message:  fmt.Sprint(retryAfter),
				Code:     r.StatusCode,
			}
		}
		return nil
	default:
		return errorResponse
	}
}
