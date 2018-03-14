// Package test contains various small helper functions that are useful when
// writing tests.
package test

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ErrorContains checks if the error message in out contains the text in
// want.
//
// This is safe when out is nil. Use an empty string for want if you want to
// test that err is nil.
func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

// Read data from a file.
func Read(t *testing.T, paths ...string) []byte {
	t.Helper()

	path := filepath.Join(paths...)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %v: %v", path, err)
	}
	return file
}

// TempFile creates a new temporary file and returns the path and a clean
// function to remove it.
//
//  f, clean := TempFile("some\ndata")
//  defer clean()
func TempFile(t *testing.T, data string) (string, func()) {
	t.Helper()

	fp, err := ioutil.TempFile(os.TempDir(), "gotest")
	if err != nil {
		t.Fatalf("test.TempFile: could not create file in %v: %v", os.TempDir(), err)
	}

	defer func() {
		err := fp.Close()
		if err != nil {
			t.Fatalf("test.TempFile: close: %v", err)
		}
	}()

	_, err = fp.WriteString(data)
	if err != nil {
		t.Fatalf("test.TempFile: write: %v", err)
	}

	return fp.Name(), func() {
		err := os.Remove(fp.Name())
		if err != nil {
			t.Errorf("test.TempFile: cannot remove %#v: %v", fp.Name(), err)
		}
	}
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
func HTTP(t *testing.T, req *http.Request, h http.Handler) *httptest.ResponseRecorder {
	t.Helper()

	rr := httptest.NewRecorder()
	if req == nil {
		var err error
		req, err = http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("cannot make request: %v", err)
		}
	}

	h.ServeHTTP(rr, req)
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
