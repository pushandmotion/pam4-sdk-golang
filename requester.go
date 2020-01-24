package pam4sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/3dsinteractive/gorequest"
)

// IRequester interface for http request
type IRequester interface {
	Get(path string, params map[string]string) (string, error)
	GetR(path string, params map[string]string) (*http.Response, string, error)
	GetRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error)
	GetRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error)
	Post(path string, params map[string]string) (string, error)
	PostR(path string, params map[string]string) (*http.Response, string, error)
	PostRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error)
	PostRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error)
	PostRaw(path string, data interface{}, headers map[string]string) (*http.Response, string, error)
	PostJSON(path string, body interface{}) (string, error)
	PostJSONR(path string, body interface{}) (*http.Response, string, error)
	PostJSONRH(path string, body interface{}, headers map[string]string) (*http.Response, string, error)
	PostJSONRHC(path string, body interface{}, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error)
	PostFile(path string, filePath string, postParam string, extraData string) (string, error)
	PutJSON(path string, body interface{}) (string, error)
	PutJSONRH(path string, body interface{}, headers map[string]string) (*http.Response, string, error)
	PutJSONRHC(path string, body interface{}, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error)
	Delete(path string, params map[string]string) (string, error)
	DeleteR(path string, params map[string]string) (*http.Response, string, error)
	DeleteRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error)
	DeleteRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error)
}

// IRequesterConfig config for requester
type IRequesterConfig interface {
	Endpoint() string
	AppIDHeaderKey() string
	SecretHeaderKey() string
	AppID() string
	Secret() string
	Timeout() time.Duration
}

// CustomRequesterConfig is custom requester config
type CustomRequesterConfig struct {
	endpoint        string
	appIDHeaderKey  string
	secretHeaderKey string
	appID           string
	secret          string
	timeout         time.Duration
}

// NewCustomRequesterConfig return new custom requester config
func NewCustomRequesterConfig(
	endpoint string,
	appIDHeaderKey string,
	secretHeaderKey string,
	appID string,
	secret string,
	timeout time.Duration) *CustomRequesterConfig {
	return &CustomRequesterConfig{
		endpoint:        endpoint,
		appIDHeaderKey:  appIDHeaderKey,
		secretHeaderKey: secretHeaderKey,
		appID:           appID,
		secret:          secret,
		timeout:         timeout,
	}
}

// Endpoint return endpoint
func (cfg *CustomRequesterConfig) Endpoint() string {
	return cfg.endpoint
}

// AppIDHeaderKey return app id header key
func (cfg *CustomRequesterConfig) AppIDHeaderKey() string {
	return cfg.appIDHeaderKey
}

// SecretHeaderKey return secret header key
func (cfg *CustomRequesterConfig) SecretHeaderKey() string {
	return cfg.secretHeaderKey
}

// AppID return app id
func (cfg *CustomRequesterConfig) AppID() string {
	return cfg.appID
}

// Secret return secret
func (cfg *CustomRequesterConfig) Secret() string {
	return cfg.secret
}

// Timeout return timeout
func (cfg *CustomRequesterConfig) Timeout() time.Duration {
	return cfg.timeout
}

// Requester struct implement IRequester
type Requester struct {
	config IRequesterConfig
	logger ILogger
	req    *gorequest.SuperAgent
}

// NewRequester return new Requester
func NewRequester(config IRequesterConfig, logger ILogger) *Requester {
	return &Requester{
		config: config,
		logger: logger,
	}
}

func (rqt *Requester) setupCredential(r *gorequest.SuperAgent) *gorequest.SuperAgent {
	appIDKey, appID := rqt.config.AppIDHeaderKey(), rqt.config.AppID()
	if len(appIDKey) > 0 && len(appID) > 0 {
		r.Set(appIDKey, appID)
	}
	secretKey, secret := rqt.config.SecretHeaderKey(), rqt.config.Secret()
	if len(secretKey) > 0 && len(secret) > 0 {
		r.Set(secretKey, secret)
	}
	return r
}

func (rqt *Requester) cloneR() *gorequest.SuperAgent {
	r := rqt.req
	if r == nil {
		r = gorequest.New()
		rqt.req = r
	}
	// Timeout is relative to time.Now so we need to set every time
	r.Timeout(rqt.config.Timeout())
	return r.Clone()
}

// Get make a GET request
func (rqt *Requester) Get(path string, params map[string]string) (string, error) {
	_, body, err := rqt.GetR(path, params)
	return body, err
}

func (rqt *Requester) GetR(path string, params map[string]string) (*http.Response, string, error) {
	res, body, err := rqt.GetRH(path, params, nil)
	return res, body, err
}

func (rqt *Requester) GetRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error) {
	return rqt.GetRHC(path, params, headers, nil)
}

func (rqt *Requester) GetRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	r := rqt.cloneR()

	url := fmt.Sprint(rqt.config.Endpoint(), path)
	rqt.logger.Debug("[RQT GET]: " + url)
	r = r.Get(url)
	rqt.setupCredential(r)
	if params != nil {
		for key, value := range params {
			r = r.Param(key, value)
		}
	}

	if headers != nil {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}

	if len(cookies) > 0 {
		for _, c := range cookies {
			r.AddCookie(c)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return res, "", NewErrorE(rqt.logger, errs[0])
	}
	rqt.logger.Debug(fmt.Sprintf("[RQT GET-RESP]: %s %s", url, rqt.truncateLogBody(body)))
	if res.StatusCode >= 400 {
		return res, body, NewErrM(res.Status)
	}

	return res, body, nil
}

// Post make a POST request
func (rqt *Requester) Post(path string, params map[string]string) (string, error) {
	_, body, err := rqt.PostR(path, params)
	return body, err
}

func (rqt *Requester) PostR(path string, params map[string]string) (*http.Response, string, error) {
	return rqt.PostRH(path, params, nil)
}

func (rqt *Requester) PostRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error) {
	return rqt.PostRHC(path, params, headers, nil)
}

func (rqt *Requester) PostRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	r := rqt.cloneR()

	u := fmt.Sprint(rqt.config.Endpoint(), path)
	rqt.logger.Debug("[RQT POST]: " + u)
	r = r.Post(u)
	rqt.setupCredential(r)
	if params != nil {
		postData := url.Values{}
		for key, value := range params {
			postData.Add(key, value)
		}
		r = r.Send(postData.Encode())
	}

	if headers != nil {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}

	if len(cookies) > 0 {
		for _, c := range cookies {
			r.AddCookie(c)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return res, "", NewErrorE(rqt.logger, errs[0])
	}
	rqt.logger.Debug(fmt.Sprintf("[RQT POST-RESP]: %s %s", u, rqt.truncateLogBody(body)))
	if res.StatusCode >= 400 {
		return res, body, NewErrM(res.Status)
	}
	return res, body, nil
}

func (rqt *Requester) PostRaw(path string, data interface{}, headers map[string]string) (*http.Response, string, error) {
	r := rqt.cloneR()

	u := fmt.Sprint(rqt.config.Endpoint(), path)
	rqt.logger.Debug("[RQT POST]: " + u)
	r = r.Post(u)
	rqt.setupCredential(r)
	r = r.Send(data)

	if headers != nil {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return res, "", NewErrorE(rqt.logger, errs[0])
	}
	rqt.logger.Debug(fmt.Sprintf("[RQT POST-RESP]: %s %s", u, rqt.truncateLogBody(body)))
	if res.StatusCode >= 400 {
		return res, body, NewErrM(res.Status)
	}
	return res, body, nil
}

func (rqt *Requester) PostJSONRHC(path string, jsonBody interface{}, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	r := rqt.cloneR()

	url := fmt.Sprint(rqt.config.Endpoint(), path)
	rqt.logger.Debug("[RQT POST]: " + url)
	r = r.Post(url)
	rqt.setupCredential(r)
	if jsonBody != nil {
		r = r.Send(jsonBody)
	}

	if headers != nil {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}

	if len(cookies) > 0 {
		for _, c := range cookies {
			r.AddCookie(c)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		for _, err := range errs {
			rqt.logger.Debug(fmt.Sprintf("[RQT POST-ERR]: %s", err.Error()))
		}
		if res != nil {
			return res, body, NewErrM(res.Status)
		} else {
			return res, body, NewErrM(errs[0].Error())
		}
	}

	rqt.logger.Debug(fmt.Sprintf("[RQT POST-RESP]: %s %s", url, rqt.truncateLogBody(body)))

	if res.StatusCode >= 400 {
		return res, body, NewErrM(res.Status)
	}
	return res, body, nil
}

// PostJSONRH make a POST request with JSON body
func (rqt *Requester) PostJSONRH(path string, jsonBody interface{}, headers map[string]string) (*http.Response, string, error) {
	return rqt.PostJSONRHC(path, jsonBody, headers, nil)
}

func (rqt *Requester) PostJSONR(path string, jsonBody interface{}) (*http.Response, string, error) {
	return rqt.PostJSONRH(path, jsonBody, nil)
}

// PostJSON make a POST request with JSON body
func (rqt *Requester) PostJSON(path string, jsonBody interface{}) (string, error) {
	_, body, err := rqt.PostJSONRH(path, jsonBody, nil)
	return body, err
}

// PostFile send file using HTTP POST
func (rqt *Requester) PostFile(path string, filePath string, postParam string, extraData string) (string, error) {
	r := rqt.cloneR()

	f, err := filepath.Abs(filePath)
	if err != nil {
		return "", NewErrorE(rqt.logger, err)
	}
	bytesOfFile, err := ioutil.ReadFile(f)
	if err != nil {
		return "", NewErrorE(rqt.logger, err)
	}

	url := fmt.Sprint(rqt.config.Endpoint(), path)
	rqt.logger.Debug("[RQT POSTFILE]: " + url + " : FILEPATH:" + filePath)

	r = r.Post(url)
	rqt.setupCredential(r)
	r = r.Type("multipart")

	if extraData != "" {
		r.Send(extraData)
	}

	r.SendFile(bytesOfFile, filepath.Base(filePath), postParam)
	res, body, errs := r.End()
	if len(errs) > 0 {
		for _, e := range errs {
			err = NewErrorE(rqt.logger, e)
		}
		return "", err
	}
	rqt.logger.Debug(fmt.Sprintf("[RQT POST-RESP]: %s %s", url, rqt.truncateLogBody(body)))
	if res.StatusCode >= 400 {
		return body, NewErrM(res.Status)
	}
	return body, nil
}

// PutJSON make a PUT request with JSON body
func (rqt *Requester) PutJSON(path string, jsonBody interface{}) (string, error) {
	_, body, err := rqt.PutJSONRH(path, jsonBody, nil)
	return body, err
}

// PutJSONRH make a PUT request with JSON body
func (rqt *Requester) PutJSONRH(path string, jsonBody interface{}, headers map[string]string) (*http.Response, string, error) {
	return rqt.PutJSONRHC(path, jsonBody, headers, nil)
}

func (rqt *Requester) PutJSONRHC(path string, jsonBody interface{}, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	r := rqt.cloneR()

	url := fmt.Sprint(rqt.config.Endpoint(), path)
	rqt.logger.Debug("[RQT PUT]: " + url)
	r = r.Put(url)
	rqt.setupCredential(r)
	if jsonBody != nil {
		r = r.Send(jsonBody)
	}

	if headers != nil {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}

	if len(cookies) > 0 {
		for _, c := range cookies {
			r.AddCookie(c)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		for _, err := range errs {
			rqt.logger.Debug(fmt.Sprintf("[RQT PUT-ERR]: %s", err.Error()))
		}
		if res != nil {
			return res, body, NewErrM(res.Status)
		} else {
			return res, body, NewErrM(errs[0].Error())
		}
	}

	rqt.logger.Debug(fmt.Sprintf("[RQT PUT-RESP]: %s %s", url, rqt.truncateLogBody(body)))

	if res.StatusCode >= 400 {
		return res, body, NewErrM(res.Status)
	}
	return res, body, nil
}

// Delete make a DELETE request
func (rqt *Requester) Delete(path string, params map[string]string) (string, error) {
	_, body, err := rqt.DeleteR(path, params)
	return body, err
}

// DeleteR make a DELETE request and response http.Response
func (rqt *Requester) DeleteR(path string, params map[string]string) (*http.Response, string, error) {
	res, body, err := rqt.DeleteRH(path, params, nil)
	return res, body, err
}

// DeleteRH make a DELETE request with headers and response http.Response
func (rqt *Requester) DeleteRH(path string, params map[string]string, headers map[string]string) (*http.Response, string, error) {
	return rqt.DeleteRHC(path, params, headers, nil)
}

// DeleteRHC make a DELETE request with headers and cookies and response http.Response
func (rqt *Requester) DeleteRHC(path string, params map[string]string, headers map[string]string, cookies []*http.Cookie) (*http.Response, string, error) {
	r := rqt.cloneR()

	url := fmt.Sprint(rqt.config.Endpoint(), path)
	rqt.logger.Debug("[RQT DELETE]: " + url)
	r = r.Delete(url)
	rqt.setupCredential(r)
	if params != nil {
		for key, value := range params {
			r = r.Param(key, value)
		}
	}

	if headers != nil {
		for key, value := range headers {
			r.Header.Add(key, value)
		}
	}

	if len(cookies) > 0 {
		for _, c := range cookies {
			r.AddCookie(c)
		}
	}

	res, body, errs := r.End()
	if len(errs) > 0 {
		return res, "", NewErrorE(rqt.logger, errs[0])
	}
	rqt.logger.Debug(fmt.Sprintf("[RQT DELETE-RESP]: %s %s", url, rqt.truncateLogBody(body)))
	if res.StatusCode >= 400 {
		return res, body, NewErrM(res.Status)
	}

	return res, body, nil
}

func (rqt *Requester) truncateLogBody(body string) string {
	length := 1000
	if len(body) > length {
		return fmt.Sprintf("%s...[resp-truncated]", body[:length-1])
	}
	return body
}
