package locationiq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// library version and default headers
const (
	libraryVersion   = "v0.0.1"
	defaultUserAgent = "location-go/" + libraryVersion
	defaultLanguage  = "en"
)

// environments
const (
	Dev  = "dev"
	Prod = "prod"
)

// BackendService services
const (
	BackendService = "backend"
)

var baseURLs = map[string]map[string]string{
	Dev: {
		BackendService: "https://eu1.locationiq.com",
	},
	Prod: {
		BackendService: "https://eu1.locationiq.com",
	},
}

// struct tags and options for requests
const (
	queryTag        = "query"
	formValueTag    = "form-value"
	formFileTag     = "form-file"
	filenameTag     = "filename"
	omitemptySuffix = ",omitempty"
)

// Client defines behaviors of tripadvisor client
type Client interface {
	Search(ctx context.Context, req SearchRequest) ([]SearchResponse, error)
}

// HTTPClient defines behaviors of http client and it is useful for mocking http client for tests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Request defines behaviors of request
type Request interface{}

// client implements Client interface
type client struct {
	Language  string
	UserAgent string
	ApiKey    string
	baseURLs  map[string]string
	c         HTTPClient
}

// compile-time proof of interface implementation
var _ Client = (*client)(nil)

// NewClient creates and returns tripadvisor client
func NewClient(environment string, c HTTPClient) Client {
	cli := client{
		Language:  defaultLanguage,
		UserAgent: defaultUserAgent,
		baseURLs:  baseURLs[environment],
		c:         c,
	}

	if cli.c == nil {
		cli.c = http.DefaultClient
	}

	return &cli
}

// request does every method's http request and returns error and http status code
func (c *client) request(ctx context.Context, method string, url string, header map[string]string, req Request, res interface{}) (int, error) {
	// get form values and form files
	formValues := getFormValues(req)
	formFiles := getFormFiles(req)
	requestHasFormData := len(formValues) > 0 || len(formFiles) > 0

	// create payload with form data or json
	payload := bytes.Buffer{}
	if requestHasFormData {
		w := multipart.NewWriter(&payload)

		for tag, value := range formValues {
			_ = w.WriteField(tag, value)
		}

		for tag, formFile := range formFiles {
			f, err := w.CreateFormFile(tag, formFile.filename)
			if err != nil {
				return 0, fmt.Errorf("creating form file failed, %s", err.Error())
			}
			r := bytes.NewBuffer(formFile.bytes)
			_, err = io.Copy(f, r)
			if err != nil {
				return 0, fmt.Errorf("writing to form file failed, %s", err.Error())
			}
		}

		_ = w.Close()

		header["Content-Type"] = w.FormDataContentType()
	} else {
		err := json.NewEncoder(&payload).Encode(req)
		if err != nil {
			return 0, fmt.Errorf("decoding request failed, %s", err.Error())
		}

		header["Content-Type"] = "application/json"
	}

	// create http request
	httpReq, err := http.NewRequestWithContext(ctx, method, url, &payload)
	if err != nil {
		return 0, fmt.Errorf("creating http request failed, %s", err.Error())
	}

	// add default header
	httpReq.Header.Set("Accept-Language", c.Language)
	httpReq.Header.Set("User-Agent", c.UserAgent)

	// add header
	for k, v := range header {
		httpReq.Header.Set(k, v)
	}

	// add query data
	q := httpReq.URL.Query()

	qd := getQueryValues(req)
	for k, v := range qd {
		q.Add(k, v)
	}

	httpReq.URL.RawQuery = q.Encode()

	// do http request
	httpRes, err := c.c.Do(httpReq)

	if err != nil {
		return 0, fmt.Errorf("doing http request failed, %s", err.Error())
	}
	defer func() {
		_ = httpRes.Body.Close()
	}()
	httpStatusCode := httpRes.StatusCode

	// decode http response
	err = json.NewDecoder(httpRes.Body).Decode(res)
	if err != nil {
		return httpStatusCode, fmt.Errorf("decoding http response failed, %s", err.Error())
	}

	// return http status errors
	if httpStatusCode >= 400 {
		return httpStatusCode, fmt.Errorf("http request failed with status: %s status code: %d", httpRes.Status, httpRes.StatusCode)
	}

	return httpStatusCode, nil
}

func getFormValues(req interface{}) map[string]string {
	return getValues(formValueTag, req)
}

func getQueryValues(req interface{}) map[string]string {
	return getValues(queryTag, req)
}

func getValues(tagName string, req interface{}) map[string]string {
	res := map[string]string{}

	e := reflect.ValueOf(req).Elem()

	for i := 0; i < e.NumField(); i++ {
		tf := e.Type().Field(i)

		t := tf.Tag.Get(tagName)
		if t == "" || t == "-" {
			continue
		}

		var v string
		var isEmpty bool
		vf := e.Field(i)
		switch vf.Interface().(type) {
		case int, int32, int64:
			iv := vf.Int()
			v = strconv.FormatInt(iv, 10)

			isEmpty = iv == 0
		case uint, uint32, uint64:
			uiv := vf.Uint()
			v = strconv.FormatUint(uiv, 10)

			isEmpty = uiv == 0
		case float32:
			fv := vf.Float()
			v = strconv.FormatFloat(fv, 'f', 6, 32)

			isEmpty = fv == 0
		case float64:
			fv := vf.Float()
			v = strconv.FormatFloat(vf.Float(), 'f', 6, 64)

			isEmpty = fv == 0
		case []byte:
			bv := vf.Bytes()
			v = string(bv)

			isEmpty = len(bv) == 0
		case string:
			v = vf.String()

			isEmpty = len(v) == 0
		}

		if strings.HasSuffix(t, omitemptySuffix) {
			if isEmpty {
				continue
			}

			t = strings.TrimSuffix(t, omitemptySuffix)
		}

		res[t] = v
	}

	return res
}

type formFile struct {
	filename string
	bytes    []byte
}

var typeOfBytes = reflect.TypeOf([]byte(nil))

func getFormFiles(req interface{}) map[string]formFile {
	res := map[string]formFile{}

	e := reflect.ValueOf(req).Elem()

	for i := 0; i < e.NumField(); i++ {
		tf := e.Type().Field(i)

		t := tf.Tag.Get(formFileTag)
		if t == "" || t == "-" {
			continue
		}

		vf := e.Field(i)
		if vf.Kind() != reflect.Slice && vf.Type() != typeOfBytes {
			continue
		}

		filename := tf.Tag.Get(filenameTag)
		if filename == "" {
			continue
		}

		res[t] = formFile{
			filename: filename,
			bytes:    vf.Bytes(),
		}
	}

	return res
}
