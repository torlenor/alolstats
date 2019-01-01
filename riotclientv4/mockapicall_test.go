package riotclientv4

var apiCallReturnJSON []byte
var apiCallReturnErr error

var lastAPICallPath string
var lastAPICallMethod string
var lastAPICallBody string

func (c *RiotClientV4) mockAPICall(path string, method string, body string) (r []byte, e error) {
	lastAPICallPath = path
	lastAPICallMethod = method
	lastAPICallBody = body

	return apiCallReturnJSON, apiCallReturnErr
}
