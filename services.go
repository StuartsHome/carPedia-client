package carpedia

import (
	"context"
	"fmt"
	"net/http"
)

type AppService service

type App struct {
	Name string
}

func (s *AppService) Get(ctx context.Context, uri string) (*http.Response, error) {
	u := "desc"
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return s.client.client.Do(req)
}

func (s *AppService) GetById(ctx context.Context, id int) (*http.Response, error) {
	u := fmt.Sprintf("desc/%d", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return s.client.client.Do(req)

}
