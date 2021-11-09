package carpedia

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

/*
	FPrintf writes to the io.Writer instance passed in
	Printf writes to the standard output
*/

const (
	baseURLPath = "/api-test"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string) {
	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	// apiHandler := http.NewServeMux()
	opts := ClientOpts{}
	client = NewClient(opts)

	url, _ := url.Parse(server.URL + "/")
	client.opts.BaseURL = url

	return client, mux, server.URL

}

// func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
// 	// mux = http.NewServeMux()

// 	// apiHandler := http.NewServeMux()
// 	// apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
// 	// apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
// 	// 	fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:") // Fprint writes to a io.Writer instance
// 	// 	fmt.Fprintln(os.Stderr)
// 	// })

// 	// // server is a test HTTP server used to provide mock API responses
// 	// server := httptest.NewServer(apiHandler)
// 	// opts := ClientOpts{}
// 	// client = NewClient(opts)
// 	// url, _ := url.Parse(server.URL + baseURLPath + "/")
// 	// client.opts.BaseURL = url

// 	// return client, mux, server.URL, server.Close
// 	return client, mux, "", nil

// }
