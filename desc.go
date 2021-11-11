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

func (d *DescService) GetDesc(ctx context.Context, opts *DescOptions) (*http.Response, error) {
	u := "desc"
	req, err := d.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return d.client.client.Do(req)

}

func (d *DescService) AddDesc(ctx context.Context, input *Desc, opts *DescOptions) (*http.Response, error) {

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
	b, err := d.client.client.Do(req)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (d *DescService) AddDescs(ctx context.Context, input *[]Desc, opts *DescOptions) ([]*http.Response, []error) {

	var statuses []*http.Response
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

func (d *DescService) GetDescById(ctx context.Context, id int, opts *DescOptions) (*http.Response, error) {

	u := fmt.Sprintf("desc/%d", id)

	req, err := d.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return d.client.client.Do(req)

}
