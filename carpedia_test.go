package carpedia

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

/*
	FPrintf writes to the io.Writer instance passed in
	Printf writes to the standard output
*/

const (
	baseURLPath = "/home"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func(), opts ClientOpts) {
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:") // Fprint writes to a io.Writer instance
		fmt.Fprintln(os.Stderr)
	})

	// server is a test HTTP server used to provide mock API responses
	server := httptest.NewServer(apiHandler)

	client = NewClient(nil, opts)
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.opts.BaseURL = url

	return client, mux, server.URL, server.Close, opts

}
