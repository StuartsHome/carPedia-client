package carpedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*

	Httptest NewServer 		- starts and returns a new server (pass in a handler)

*/

/*
func setup() (client *Client, mux *http.ServeMux, serverURL string) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {})
	server := httptest.NewServer(apiHandler)

	opts := ClientOpts{}
	client = NewClient(opts)

	url, _ := url.Parse(server.URL + "/")
	client.opts.BaseURL = url

	return client, mux, server.URL
}
*/

func TestCarUseSetup(t *testing.T) {
	client, mux, _ := setup()
	body := "tony"

	mux.HandleFunc("/car", func(w http.ResponseWriter, r *http.Request) {
		// r.Body = ioutil.NopCloser(bytes.NewBufferString(body))
		w.Header().Add("Content-Type", "application/json")
		// w.Header().Set("text/plain", "charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(body))
		fmt.Fprint(w, "hello")
		testMethod(t, r, http.MethodGet)
	})

	ctx := context.Background()
	opts := CarOptions{}
	cars, err := client.Car.Car(ctx, &opts)
	require.Nil(t, err)

	// body, err := ioutil.ReadAll(cars.Body) -> read body into bytes

	// var tempData interface{}
	// json.NewDecoder(cars.Body).Decode(&tempData)
	aa, _ := io.ReadAll(cars.Body)
	fmt.Println(string(aa))

}

func TestCar(t *testing.T) {

	// This allows each test to create its own handler by changing handler variable
	handler := http.NotFound
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer server.Close()

	opts := ClientOpts{}
	client := NewClient(opts)

	handler = func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" {
			t.Error("Bad path!")
		}
		io.WriteString(w, `{"car": "vw"}`)
	}

	car, err := client.Car.Car(context.TODO(), nil)
	if err != nil {
		t.Errorf("error calling Car method")
	}

	assert.Equal(t, car.StatusCode, 200)
	require.NotNil(t, car)

	var tempData interface{}
	json.NewDecoder(car.Body).Decode(&tempData)
	fmt.Println(tempData)

}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("header.Get(%q) returned %q, want %q", header, got, want)
	}
}
