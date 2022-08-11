package requests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var _ ClientAPI = (*Client)(nil)

type (
	// http client interface
	ClientAPI interface {
		Get(context.Context, string, map[string]string) ([]byte, error)
	}

	// http client struct
	Client struct {
		h *http.Client
	}

	errResponse struct {
		Message string
	}
)

//NewClient initialize new Http client
func NewClient() (*Client, error) {
	h := &http.Client{
		// set default 60 seconds timeout
		Timeout: time.Minute,
	}
	return &Client{
		h: h,
	}, nil
}

func (c *Client) Get(ctx context.Context, url string, params map[string]string) ([]byte, error) {

	// create new http request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(ctx, "error", "error while creating new http request", err.Error())
		return nil, err
	}
	// fill the params
	values := req.URL.Query()
	for key, element := range params {
		values.Add(key, element)
	}
	req.URL.RawQuery = values.Encode()

	//make request
	res, err := c.h.Do(req)
	if err != nil {
		log.Println(ctx, "error", "error while making http GET request", err.Error())
		return nil, err
	}

	defer res.Body.Close()

	//read bytes from response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(ctx, "error", "error while reading response body", err.Error())
		return nil, err
	}

	if res.StatusCode != 200 {
		// got error from server
		return nil, unwrapHttpError(ctx, body)
	}

	return body, nil
}

//unwrap the error message
func unwrapHttpError(ctx context.Context, body []byte) error {
	var errRes errResponse
	err := json.Unmarshal(body, &errRes)
	if err != nil {
		log.Println(ctx, "error", "error while un marshalling errored response body", err.Error(), "body", string(body))
		return fmt.Errorf("error while un marshalling errored response body, err: %v", err.Error())
	}
	return fmt.Errorf("error message: %v", errRes.Message)
}
