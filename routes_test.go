package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAccountPayload(t *testing.T) {
	// Wil test the unmarshalling of the account into the route payload
	// https://gosamples.dev/struct-to-io-reader/
	byt, _ := json.Marshal(&UserAccount{})
	reader := bytes.NewReader(byt)
	ctx := &gin.Context{
		Request: &http.Request{
			Method: "POST",
			URL:    &url.URL{Path: "some/random/texzt/url"},
			Proto:  "HTTP/1.1",
			Header: map[string][]string{
				"Accept-Encoding": {"application/json"},
			},
			Body: io.NopCloser(reader),
		},
	}
	AccountPayload(ctx)
}
