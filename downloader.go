package downloader

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

type Downloader struct {
	UserAgent         string
	Async             bool
	ID                uint32
	client            *http.Client
	wg                *sync.WaitGroup
	lock              *sync.RWMutex
	requestCallbacks  []RequestCallback
	responseCallbacks []ResponseCallback
	errorCallbacks    []ErrorCallback
}

// RequestCallback is a type alias for OnRequest callback functions
type RequestCallback func(*http.Request)

// ResponseCallback is a type alias for OnResponse callback functions
type ResponseCallback func(*http.Response)

// ErrorCallback is a type alias for OnError callback functions
type ErrorCallback func(*http.Response, error)

// ProxyFunc is a type alias for proxy setter functions.
type ProxyFunc func(*http.Request) (*url.URL, error)

var downloaderCounter uint32

func NewDownloader(options ...func(*Downloader)) *Downloader {
	d := &Downloader{}
	d.Init()
	for _, f := range options {
		f(d)
	}
	return d
}

func (dl *Downloader) Init() {
	dl.UserAgent = ""
	dl.client = &http.Client{}
	dl.wg = &sync.WaitGroup{}
	dl.lock = &sync.RWMutex{}
	dl.ID = atomic.AddUint32(&downloaderCounter, 1)
}

func UserAgent(us string) func(*Downloader) {
	return func(dl *Downloader) {
		dl.UserAgent = us
	}
}

func Async(async bool) func(*Downloader) {
	return func(dl *Downloader) {
		dl.Async = async
	}
}

func (dl *Downloader) Wait() {
	dl.wg.Wait()
}

func (dl *Downloader) Get(URL string) error {
	return dl.scrape(URL, "GET", nil, nil)
}

func (dl *Downloader) POST(URL string, requestData map[string]string) error {
	return dl.scrape(URL, "POST", createFormReader(requestData), nil)
}

func (dl *Downloader) scrape(u, method string, requestData io.Reader, hdr http.Header) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}
	if hdr == nil {
		hdr = http.Header{"User-Agent": []string{dl.UserAgent}}
	}
	rc, ok := requestData.(io.ReadCloser)
	if !ok && requestData != nil {
		rc = ioutil.NopCloser(requestData)
	}
	req := &http.Request{
		Method:     method,
		URL:        parsedURL,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       rc,
		Host:       parsedURL.Host,
		Header:     map[string][]string{},
	}
	setRequestBody(req, requestData)
	u = parsedURL.String()
	dl.wg.Add(1)
	if dl.Async {
		go dl.fetch(u, method, requestData, hdr, req)
		return nil
	}
	return dl.fetch(u, method, requestData, hdr, req)
}

func (dl *Downloader) fetch(u, method string, requestData io.Reader, hdr http.Header, req *http.Request) error {
	defer dl.wg.Done()

	dl.handleOnRequest(req)

	if method == "POST" && req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "*/*")
	}

	response, err := dl.client.Do(req)
	err = dl.handlerOnError(response, err, req)
	if err != nil {
		return err
	}
	dl.handleOnResponse(response)
	return nil
}

func setRequestBody(req *http.Request, body io.Reader) {
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return ioutil.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return ioutil.NopCloser(&r), nil
			}
		}
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = http.NoBody
			req.GetBody = func() (io.ReadCloser, error) { return http.NoBody, nil }
		}
	}
}

// OnRequest
func (dl *Downloader) OnRequest(f RequestCallback) {
	dl.lock.Lock()
	if dl.requestCallbacks == nil {
		dl.requestCallbacks = make([]RequestCallback, 0, 4)
	}
	dl.requestCallbacks = append(dl.requestCallbacks, f)
	dl.lock.Unlock()
}

// OnResponse
func (dl *Downloader) OnResponse(f ResponseCallback) {
	dl.lock.Lock()
	if dl.responseCallbacks == nil {
		dl.responseCallbacks = make([]ResponseCallback, 0, 4)
	}
	dl.responseCallbacks = append(dl.responseCallbacks, f)
	dl.lock.Unlock()
}

// onError
func (dl *Downloader) OnError(f ErrorCallback) {
	dl.lock.Lock()
	if dl.errorCallbacks == nil {
		dl.errorCallbacks = make([]ErrorCallback, 0, 4)
	}
	dl.errorCallbacks = append(dl.errorCallbacks, f)
	dl.lock.Unlock()
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

func (dl *Downloader) handleOnRequest(r *http.Request) {
	for _, f := range dl.requestCallbacks {
		f(r)
	}
}

func (dl *Downloader) handleOnResponse(resp *http.Response) {
	for _, f := range dl.responseCallbacks {
		f(resp)
	}
}

func (dl *Downloader) handlerOnError(response *http.Response, err error, request *http.Request) error {
	if err == nil && response.StatusCode < 203 {
		return nil
	}
	if err == nil && response.StatusCode >= 203 {
		err = errors.New(http.StatusText(response.StatusCode))
	}
	if response == nil {
		response = &http.Response{
			Request: request,
		}
	}
	if response.Request == nil {
		response.Request = request
	}
	for _, f := range dl.errorCallbacks {
		f(response, err)
	}
	return err
}
