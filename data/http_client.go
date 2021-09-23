package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// HttpClient defines a simple interface to execute an http request and map the request and response objects.
type HttpClient interface {
	// Do will execute an http request.
	Do(ctx context.Context, method, endpoint string, expStatus int, req interface{}, out interface{}) error
}

type client struct {
	c *http.Client
}

// NewClient will setup and return a new http client.
func NewClient(c *http.Client) *client {
	return &client{
		c: c,
	}
}

// Do will execute an http request and validate the status matches expStatus.
//
// if req is empty no request body will be added, if out is empty, the response will not be mapped.
func (c *client) Do(ctx context.Context, method, endpoint string, expStatus int, req interface{}, out interface{}) error {
	rdr := &bytes.Buffer{}
	if req != nil {
		if err := json.NewEncoder(rdr).Encode(req); err != nil {
			return errors.Wrapf(err, "failed to encode request for '%s' '%s'", method, endpoint)
		}
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, endpoint, rdr)
	if err != nil {
		return errors.Wrapf(err, "failed to create http request for '%s' '%s'", method, endpoint)
	}
	httpReq.Header.Add("Content-Type", "application/json")

	resp, err := c.c.Do(httpReq)
	if err != nil {
		return errors.Wrapf(err, "failed to send request to for '%s' '%s'", method, endpoint)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != expStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error for '%s' '%s'. Status Received : '%d', Status Expected : '%d'. \nBody: %s", method, endpoint, resp.StatusCode, expStatus, body)
	}
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return errors.Wrapf(err, "failed to decode response for '%s' '%s'", method, endpoint)
		}
	}
	return nil
}
