package carpedia

import (
	"context"
	"net/http"
)

type DescService service

type DescOptions struct {
}

func (d *DescService) Desc(ctx context.Context, opts *DescOptions) (*http.Response, error) {
	u := "desc"
	req, err := d.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return d.client.client.Do(req)

}

// TODO:
func (d *DescService) AddDesc(ctx context.Context, opts *DescOptions) (*http.Response, error) {
	return nil, nil
}
