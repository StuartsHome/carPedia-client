package carpedia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DescService service

type DescOptions struct {
}

func (d *DescService) GetDesc(ctx context.Context, opts *DescOptions) (*Response, error) {
	u := "desc"
	req, err := d.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return d.client.defaultDo(ctx, req)

}

func (d *DescService) AddDesc(ctx context.Context, input *Desc, opts *DescOptions) (*Response, error) {

	u := "desc"
	body, err := json.Marshal(&input)
	if err != nil {
		return nil, err
	}

	req, err := d.client.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	b, err := d.client.defaultDo(ctx, req)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (d *DescService) AddDescs(ctx context.Context, input *[]Desc, opts *DescOptions) ([]*Response, []error) {

	var statuses []*Response
	var errs []error
	for _, desc := range *input {
		response, err := d.AddDesc(ctx, &desc, opts)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		statuses = append(statuses, response)
	}
	if len(errs) != 0 {
		return nil, errs
	}
	return statuses, nil
}

func (d *DescService) GetDescById(ctx context.Context, id int, opts *DescOptions) (*Response, error) {

	u := fmt.Sprintf("desc/%d", id)

	req, err := d.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return d.client.defaultDo(ctx, req)

}
