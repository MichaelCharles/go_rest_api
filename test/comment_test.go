// +build e2e

package test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetComments(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get("http://localhost:8080/api/health")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 200, resp.StatusCode())
}

func TestPostComment(t *testing.T) {
	client := resty.New()
	resp, err := client.R().
		SetBody(`{"slug": "/test-comment", "author": "12345", "body": "this is a pen"}`).
		Post(BASE_URL + "/api/comment")

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode())
}
