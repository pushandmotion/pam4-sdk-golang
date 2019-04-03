package pam4sdk

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/3dsinteractive/jason"
	"github.com/3dsinteractive/testify/assert"
	"github.com/3dsinteractive/testify/suite"
)

type RequesterTestSuite struct {
	suite.Suite
}

func TestRequesterTestSuite(t *testing.T) {
	suite.Run(t, new(RequesterTestSuite))
}

func (ts *RequesterTestSuite) requesterConfig(url string) *MockRequesterConfig {
	cfg := NewMockRequesterConfig()
	cfg.On("AppIDHeaderKey").Return("x-app-id")
	cfg.On("SecretHeaderKey").Return("x-secret")
	cfg.On("AppID").Return("my-app-id-1234")
	cfg.On("Secret").Return("my-secret-1234")
	cfg.On("Endpoint").Return(url)
	cfg.On("Timeout").Return(time.Second * 2)
	return cfg
}

func (ts *RequesterTestSuite) TestGET_GivenNoParams_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	result, err := rqt.Get("/abc", nil)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}

func (ts *RequesterTestSuite) TestGET_GivenParams_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc?a=abc&b=def", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))
		is.Equal("abc", req.URL.Query().Get("a"))
		is.Equal("def", req.URL.Query().Get("b"))

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	params := map[string]string{
		"a": "abc",
		"b": "def",
	}
	result, err := rqt.Get("/abc", params)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}

func (ts *RequesterTestSuite) TestPOST_GivenNoParams_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	result, err := rqt.Post("/abc", nil)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}

func (ts *RequesterTestSuite) TestPOST_GivenParams_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))

		err := req.ParseForm()
		if is.NoError(err) {
			is.Equal("abc", req.Form.Get("a"))
			is.Equal("def", req.Form.Get("b"))
		}

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	params := map[string]string{
		"a": "abc",
		"b": "def",
	}
	result, err := rqt.Post("/abc", params)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}

func (ts *RequesterTestSuite) TestPostJSON_GivenNoParams_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	result, err := rqt.PostJSON("/abc", nil)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}

func (ts *RequesterTestSuite) TestPostJSON_GivenParams_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))

		body, err := ioutil.ReadAll(req.Body)
		v, err := jason.NewObjectFromBytes(body)
		if is.NoError(err) {
			a, _ := v.GetString("a")
			is.Equal("abc", a)
			b, _ := v.GetString("b")
			is.Equal("def", b)
		}

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	params := map[string]string{
		"a": "abc",
		"b": "def",
	}
	result, err := rqt.PostJSON("/abc", params)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}

func (ts *RequesterTestSuite) TestPostJSON_GivenParamsAndHeaders_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))

		is.Equal("ccc", req.Header.Get("c"))
		is.Equal("ddd", req.Header.Get("d"))

		body, err := ioutil.ReadAll(req.Body)
		v, err := jason.NewObjectFromBytes(body)
		if is.NoError(err) {
			a, _ := v.GetString("a")
			is.Equal("abc", a)
			b, _ := v.GetString("b")
			is.Equal("def", b)
		}

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	params := map[string]string{
		"a": "abc",
		"b": "def",
	}
	headers := map[string]string{
		"c": "ccc",
		"d": "ddd",
	}
	_, result, err := rqt.PostJSONRH("/abc", params, headers)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}

func (ts *RequesterTestSuite) TestPostJSON_GivenParamsAndHeadersAndCookies_ExpectCorrectResult() {
	is := assert.New(ts.T())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		is.Equal("/abc", req.URL.String())
		is.Equal("my-app-id-1234", req.Header.Get("x-app-id"))
		is.Equal("my-secret-1234", req.Header.Get("x-secret"))

		is.Equal("ccc", req.Header.Get("c"))
		is.Equal("ddd", req.Header.Get("d"))

		c, _ := req.Cookie("cookie1")
		is.Equal("cookie_value1", c.Value)
		c, _ = req.Cookie("cookie2")
		is.Equal("cookie_value2", c.Value)

		body, err := ioutil.ReadAll(req.Body)
		v, err := jason.NewObjectFromBytes(body)
		if is.NoError(err) {
			a, _ := v.GetString("a")
			is.Equal("abc", a)
			b, _ := v.GetString("b")
			is.Equal("def", b)
		}

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	rqt := NewRequester(ts.requesterConfig(server.URL), NewLoggerSimple())
	params := map[string]string{
		"a": "abc",
		"b": "def",
	}
	headers := map[string]string{
		"c": "ccc",
		"d": "ddd",
	}
	cookies := []*http.Cookie{
		&http.Cookie{
			Name:  "cookie1",
			Value: "cookie_value1",
		},
		&http.Cookie{
			Name:  "cookie2",
			Value: "cookie_value2",
		},
	}
	_, result, err := rqt.PostJSONRHC("/abc", params, headers, cookies)
	if is.NoError(err) {
		is.Equal("OK", result)
	}
}
