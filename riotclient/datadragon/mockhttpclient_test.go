package riotclientdd

import "net/http"

type mockHTTPClient struct {
	cnt int16

	// What the Get request was
	getURL  string
	getURL2 string

	// What the mock will return on first request
	response *http.Response
	err      error

	// What the mock will return on second request
	response2 *http.Response
	err2      error
}

func (m *mockHTTPClient) Get(url string) (resp *http.Response, err error) {

	if m.cnt == 0 {
		m.getURL = url
		m.cnt++
		return m.response, m.err
	}
	m.getURL2 = url
	return m.response2, m.err2
}
