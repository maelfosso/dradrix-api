package requests

import (
	"io"
	"net/http"
)

type IHttpRequester interface {
	Get(url string) (*http.Response, error)
	Put(url string, contentLength int64, body io.Reader)
	Delete(url string) (*http.Response, error)
}

type HttpRequester struct{}

func (httpRequester HttpRequester) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func (httpRequester HttpRequester) Put(url string, contentLength int64, body io.Reader) (*http.Response, error) {
	putRequest, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	putRequest.ContentLength = contentLength
	return http.DefaultClient.Do(putRequest)
}

func (httpRequester HttpRequester) Delete(url string) (*http.Response, error) {
	deleteRequest, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(deleteRequest)
}
