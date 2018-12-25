package riotclientv4

var apiCallReturnJSON []byte
var apiCallReturnErr error

func (c *RiotClientV4) mockAPICall(path string, method string, body string) (r []byte, e error) {
	return apiCallReturnJSON, apiCallReturnErr
}
