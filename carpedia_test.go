package carpedia

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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

// TODO: tidy up
func testNewRequestAndDoFailure(t *testing.T, methodName string, client *Client, f func() (*http.Response, error)) {
	t.Helper()
	if methodName == "" {
		t.Error("testNewRequestAndDoFailure: method name empty")
	}

	client.opts.BaseURL.Path = ""
	resp, err := f()
	if err != nil {
		t.Errorf("client.BaseURL.Path:'' %v resp: %#v, want: nil", methodName, resp)
	}
	if err != nil {
		t.Errorf("client.BaseURL.Path:'' %v err: nil, want: error", methodName)
	}

	client.opts.BaseURL.Path = baseURLPath
	resp, err = f()
	if want := http.StatusForbidden; resp == nil || resp.StatusCode != want {
		if resp != nil {
			t.Errorf("resp = %#v, want StatusCode=%v", resp, want)
		} else {
			t.Errorf("resp = nil, want StatusCode=%v", want)
		}
	}
	if err == nil {
		t.Error("err = nil, want error", methodName)
	}
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
