package carpedia

import (
	"context"
	"net/http"
)

type CarService service

type CarOptions struct {
}

func (c *CarService) Car(ctx context.Context, opts *CarOptions) (*Response, error) {

	// params, err := qs.Values(opts)
	// if err != nil {
	// 	return nil, err
	// }
	// u := fmt.Sprint("home", params.Encode())

	u := "car"
	req, err := c.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	// return c.client.client.Do(req)
	return c.client.defaultDo(ctx, req)
}
