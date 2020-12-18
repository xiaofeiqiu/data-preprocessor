package restutils

import (
	"bytes"
	"fmt"
	"github.com/xiaofeiqiu/mlstock/lib/log"
	"io/ioutil"
	"net/http"
)

const GET = "GET"

type HttpClient struct {
	client *http.Client
}

func (hc *HttpClient) doRequest(method string, url string, extraHeaders map[string]string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		log.Error("GetRequest", fmt.Sprintf("Error creating new request: %v", err))
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip,deflate")

	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	res, err := hc.client.Do(req)
	if err != nil {
		log.Error("GetRequest", fmt.Sprintf("Error executing the reqeust: %v", err))
		return 0, nil, err
	}
	defer res.Body.Close()

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("GetRequest", fmt.Sprintf("Error reading the response body: %v", err))
		return 0, nil, err
	}

	return res.StatusCode, resBodyBytes, nil
}

func (hc *HttpClient) DoGet(url string, extraHeaders map[string]string) (int, []byte, error) {
	return hc.doRequest(http.MethodGet, url, extraHeaders, nil)
}

func (hc *HttpClient) DoPut(url string, extraHeaders map[string]string, body []byte) (int, []byte, error) {
	return hc.doRequest(http.MethodPut, url, extraHeaders, body)
}

func (hc *HttpClient) DoPost(url string, extraHeaders map[string]string, body []byte) (int, []byte, error) {
	return hc.doRequest(http.MethodPost, url, extraHeaders, body)
}

func (hc *HttpClient) DoDelete(url string, extraHeaders map[string]string,body []byte) (int, []byte, error) {
	return hc.doRequest(http.MethodDelete, url, extraHeaders, body)
}


