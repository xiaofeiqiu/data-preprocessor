package restutils

import (
	"bytes"
	"encoding/json"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"io/ioutil"
	"net/http"
)

const GET = "GET"

type HttpClient struct {
	Client *http.Client
}

func (hc *HttpClient) doRequest(method string, url string, extraHeaders map[string]string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		log.Error("GetRequest", err, "error creating new request")
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Accept", "*/*")

	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	res, err := hc.Client.Do(req)
	if err != nil {
		log.Error("GetRequest", err, "error executing the reqeust")
		return 0, nil, err
	}
	defer res.Body.Close()

	resBodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("GetRequest", err, "error reading the response body")
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

func ResponseWithJson(w http.ResponseWriter, code int, body interface{}) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshaling body to json"))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonBody)
	return
}

func Is2xxStatusCode(status int) bool {
	if status >= 200 && status < 300 {
		return true
	}
	return false
}

func Is4xxStatusCode(status int) bool {
	if status >= 400 && status < 500 {
		return true
	}
	return false
}

func Is5xxStatusCode(status int) bool {
	if status >= 500 && status < 600 {
		return true
	}
	return false
}