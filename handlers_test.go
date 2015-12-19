package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

func TestCreateMatch(t *testing.T) {
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter)))
	defer server.Close()

	body := []byte("{\n  \"gridsize\": 19,\n  \"players\": [\n    {\n      \"color\": \"white\",\n      \"name\": \"bob\"\n    },\n    {\n      \"color\": \"black\",\n      \"name\": \"alfred\"\n    }\n  ]\n}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createMatchHandler: %v", err)
	}

	defer res.Body.Close()

	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status 201, received %s", res.Status)
	}

	if res.Header["Location"] == nil {
		t.Error("Location header is not set")
	}

	fmt.Printf("Payload: %s", string(payload))
}
