// Package test contains various small helper functions that are useful when
// writing tests.
package test

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

// ErrorContains checks if the error message in got contains the text in
// expected.
//
// This is safe when got is nil. Use an empty string for expected if you want to
// test that err is nil.
func ErrorContains(got error, expected string) bool {
	if got == nil {
		return expected == ""
	}
	if expected == "" {
		return false
	}
	return strings.Contains(got.Error(), expected)
}

// Read data from a file.
func Read(t *testing.T, paths ...string) []byte {
	path := filepath.Join(paths...)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %v: %v", path, err)
	}
	return file
}

// HTTP sets up a HTTP test. A GET request will be made for you if req is nil.
//
// For example:
//
//     rr := test.HTTP(t, nil, MyHandler)
//
// Or for a POST request:
//
//     req, err := http.NewRequest("POST", "/v1/email", b)
//     if err != nil {
//     	t.Fatalf("cannot make request: %v", err)
//     }
//     req.Header.Set("Content-Type", ct)
//     rr := test.HTTP(t, req, MyHandler)
func HTTP(t *testing.T, req *http.Request, h http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h)
	if req == nil {
		var err error
		req, err = http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("cannot make request: %v", err)
		}
	}

	handler.ServeHTTP(rr, req)
	return rr
}

// MultipartForm writes the keys and values from params to a multipart form.
//
// Don't forget to set the Content-Type:
//
//   req.Header.Set("Content-Type", contentType)
func MultipartForm(params ...map[string]string) (b *bytes.Buffer, contentType string, err error) {
	b = &bytes.Buffer{}
	w := multipart.NewWriter(b)

	for k, v := range params[0] {
		field, err := w.CreateFormField(k)
		if err != nil {
			return nil, "", err
		}
		_, err = field.Write([]byte(v))
		if err != nil {
			return nil, "", err
		}
	}

	if len(params) > 1 {
		for k, v := range params[1] {
			field, err := w.CreateFormFile(k, k)
			if err != nil {
				return nil, "", err
			}
			_, err = field.Write([]byte(v))
			if err != nil {
				return nil, "", err
			}
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return b, w.FormDataContentType(), nil
}
