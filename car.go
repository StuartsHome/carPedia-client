package carpedia

import (
	"context"
	"fmt"
	"net/http"

	qs "github.com/google/go-querystring/query"
)

type CarService service

type CarOptions struct {
}

func (c *CarService) Car(ctx context.Context, opts *CarOptions) (*Response, error) {

	params, err := qs.Values(opts)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("home/%s", params.Encode())
	req, err := c.client.NewRequest(http.MethodGet, u, "")
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.client.defaultDo(ctx, req)
}
