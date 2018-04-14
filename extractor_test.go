package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestGetLandingPage(t *testing.T) {

	req, err := http.NewRequest("GET", "localhost:8087/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	rec := httptest.NewRecorder()
	getLandingPage(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}
	data := url.Values{}
	data.Set("url", "https://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest")
	req, err = http.NewRequest("PUT", "localhost:8087/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	getLandingPage(rec, req)
	res = rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404; got %v", res.Status)
	}

	//post request
	req, err = http.NewRequest("POST", "localhost:8087/", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rec = httptest.NewRecorder()
	getLandingPage(rec, req)
	res = rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

}

func TestMapper(t *testing.T) {

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "null string", value: "", want: ""},
		{name: "word with .", value: "prabesh.", want: "prabesh "},
		{name: "word with symbol", value: "prabesh߶", want: "prabesh"},
		{name: "word with number", value: "prabesh123", want: "prabesh"},
		{name: "word with number", value: "prabesh rasmey", want: "prabesh rasmey"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapper(tt.value); got != tt.want {
				t.Errorf("mapper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtrator(t *testing.T) {

	tests := []struct {
		name string
		body io.Reader
		out  chan string
	}{
		{name: "simple test", body: bytes.NewReader([]byte("its me prabesh")), out: make(chan string)},
		{name: "with embedded script and style", body: bytes.NewReader([]byte(`<!doctype html><html lang="en"> <body><div id="result"></div><script type="text/javascript">console.log("hello")</script> </body></html>`)), out: make(chan string)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extrator(tt.body, tt.out)
		})
	}
}

func TestValidateInt(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{name: "not a number", id: "hg", want: false},
		{name: "a number", id: "12", want: true},
		{name: "empty", id: "", want: false},
		{name: "negative number", id: "-1", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateInt(tt.id); got != tt.want {
				t.Errorf("validateInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateColumn(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{name: "Uppercase alphabet", val: "AA", want: true},
		{name: "Lowercase alphabet", val: "aa", want: false},
		{name: "Empty string", val: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateColumn(tt.val); got != tt.want {
				t.Errorf("validateColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}