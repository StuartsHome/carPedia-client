package carpedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDesc(t *testing.T) {
	client, mux, _ := setup()
	body := "{null, null}"

	mux.HandleFunc("/desc", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))

		testMethod(t, r, http.MethodGet)
	})

	ctx := context.Background()
	opts := DescOptions{}
	descs, err := client.Desc.GetDesc(ctx, &opts)
	require.Nil(t, err)

	actual, err := io.ReadAll(descs.Body)
	require.Nil(t, err)

	assert.Equal(t, "{null, null}", string(actual))

}

// TODO
func TestGetDescFail(t *testing.T) {

}

func TestAddDesc(t *testing.T) {
	client, mux, _ := setup()

	input := Desc{
		Id:    1,
		Title: "Test Title",
		Text:  "this is a test",
	}

	mux.HandleFunc("/desc", func(w http.ResponseWriter, r *http.Request) {
		v := Desc{Id: 1}
		json.NewDecoder(r.Body).Decode(&v)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		testMethod(t, r, http.MethodPost)

		want := Desc{Id: 1}
		if !cmp.Equal(v, want) {
			t.Errorf("body %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `{"Id": 1}`)
	})

	ctx := context.Background()
	opts := DescOptions{}

	response, err := client.Desc.AddDesc(ctx, &input, &opts)
	require.Nil(t, err)

	var got Desc
	json.NewDecoder(response.Body).Decode(&got)
	want := Desc{Id: 1}

	assert.Equal(t, want, got)
}

// TODO
func TestAddDescFail(t *testing.T) {

}

func TestAddDescs(t *testing.T) {
	client, mux, _ := setup()

	input := []Desc{
		{
			Id:    1,
			Title: "Test Title",
			Text:  "this is a test",
		},
		{
			Id:    2,
			Title: "Test Title",
			Text:  "this is a test",
		},
	}

	mux.HandleFunc("/desc", func(w http.ResponseWriter, r *http.Request) {
		v := []Desc{{Id: 1}, {Id: 2}}
		json.NewDecoder(r.Body).Decode(&v)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		testMethod(t, r, http.MethodPost)

		want := []Desc{{Id: 1}, {Id: 2}}
		if !cmp.Equal(v, want) {
			t.Errorf("body %+v, want %+v", v, want)
		}
		fmt.Fprint(w, `[{"Id": 1}, {"Id": 2}]`)
	})

	ctx := context.Background()
	opts := DescOptions{}

	responses, err := client.Desc.AddDescs(ctx, &input, &opts)
	require.Nil(t, err)

	var got []Desc
	for _, desc := range responses {
		json.NewDecoder(desc.Body).Decode(&got)
	}
	want := []Desc{
		{Id: 1},
		{Id: 2},
	}

	assert.Equal(t, want, got)

}

// TODO
func TestAddDescsFail(t *testing.T) {

}

// TODO: Add func 'testNewRequestAndDoFailure'
func TestGetDescById(t *testing.T) {
	client, mux, _ := setup()

	mux.HandleFunc("/desc/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"id":1,"title":"Test Title","text":"this is a test"}`)
	})

	ctx := context.Background()
	opts := DescOptions{}

	response, err := client.Desc.GetDescById(ctx, 1, &opts)
	require.Nil(t, err)

	var got Desc
	json.NewDecoder(response.Body).Decode(&got)
	want := Desc{Id: 1, Title: "Test Title", Text: "this is a test"}

	assert.Equal(t, want, got)

}

type Timestamp struct {
	time.Time
}

func (t Timestamp) Equal(u Timestamp) bool {
	return t.Time.Equal(u.Time)
}
