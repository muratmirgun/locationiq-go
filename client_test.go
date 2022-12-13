package locationiq

import (
	"net/http"
	"net/http/httptest"
)

var (
	testClient *client
	testMux    *http.ServeMux
	testServer *httptest.Server
)

func setupTest() {
	testClient = NewClient("dev", "", nil).(*client)
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)
	url := testServer.URL
	testClient.baseURLs = map[string]string{
		BackendService: url,
	}
}

func teardownTest() {
	testServer.Close()
}
