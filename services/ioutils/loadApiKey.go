package ioutils

import (
	"io/ioutil"
)

// load api key from file
func LoadApiKey(path string) (string, error) {
	apiKey, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(apiKey) , nil
}
