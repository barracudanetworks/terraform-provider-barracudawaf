package barracudawaf

import (
	"bytes"
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	baseURI = "restapi/v3.1" // baseURI : base endpoint for APICall
)

// BarracudaWAF : container for barracuda's WAF session state.
type BarracudaWAF struct {
	Host      string
	User      string
	Password  string
	Token     string
	UserAgent string //specifies the caller of the request
	Transport *http.Transport
}

// APIRequest : builds API request for resource.
type APIRequest struct {
	Method      string
	URL         string
	Body        interface{}
	ContentType string
}

//WAFResouceData : Container for barracuda WAF resource's data
type WAFResouceData struct {
	Token  string                            `json:"token,omitempty"`
	Object string                            `json:"object,omitempty"`
	Data   map[string]map[string]interface{} `json:"data,omitempty"`
}

// RequestError : contains information about any error from a request.
type RequestError struct {
	Message string `json:"msg,omitempty"`
	Token   string `json:"token,omitempty"`
}

// Error : returns the error message.
func (r *RequestError) Error() error {
	if r.Message != "" {
		return errors.New(r.Message)
	}

	return nil
}

// NewSession : Barracuda WAF system connection.
func NewSession(host, port, user, passwd string) *BarracudaWAF {
	var url string
	if !strings.HasPrefix(host, "http") {
		url = fmt.Sprintf("https://%s", host)
	} else {
		url = host
	}
	if port != "" {
		url = url + ":" + port
	}

	return &BarracudaWAF{
		Host:     url,
		User:     user,
		Password: passwd,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

// GetAuthToken : Check Token for Authentication Barracuda WAF APIs
func (b *BarracudaWAF) GetAuthToken() (*BarracudaWAF, error) {

	body := map[string]string{
		"username": b.User,
		"password": b.Password,
	}

	marshalJSON, err := jsonMarshal(body)

	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         "/login",
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	callRes, callErr := b.APICall(req)

	if callErr != nil {
		return nil, callErr
	}

	var resBody map[string]string
	unmarshalError := json.Unmarshal(callRes, &resBody)

	if unmarshalError != nil {
		log.Printf("[ERROR] Unable to unmarshal auth token response")
		return nil, unmarshalError
	}

	b.Token = resBody["token"] + ":"
	b.Token = "BASIC " + b64.StdEncoding.EncodeToString([]byte(b.Token))

	return b, nil
}

// CreateBarracudaWAFResource : Creates Barracuda WAF resource
func (b *BarracudaWAF) CreateBarracudaWAFResource(name string, request *APIRequest) error {
	_, err := b.postReq(request.Body, request.URL)

	return err
}

// UpdateBarracudaWAFResource : Updates Barracuda WAF resource
func (b *BarracudaWAF) UpdateBarracudaWAFResource(name string, request *APIRequest) error {
	_, err := b.putReq(request.Body, fmt.Sprintf("%s/%s", request.URL, name))

	return err
}

// UpdateBarracudaWAFSubResource : Updates Barracuda WAF sub resource
func (b *BarracudaWAF) UpdateBarracudaWAFSubResource(name string, resourceEndpoint string, request *APIRequest) error {
	_, err := b.putReq(request.Body, fmt.Sprintf("%s/%s/%s", resourceEndpoint, name, request.URL))

	return err
}

// DeleteBarracudaWAFResource : Delete Barracuda WAF resource
func (b *BarracudaWAF) DeleteBarracudaWAFResource(name string, request *APIRequest) error {
	_, err := b.deleteReq(fmt.Sprintf("%s/%s", request.URL, name))

	return err
}

// GetBarracudaWAFResource : Updates Barracuda WAF resource
func (b *BarracudaWAF) GetBarracudaWAFResource(name string, request *APIRequest) (*WAFResouceData, error) {
	data, err := b.getReq(request.URL)

	if err != nil {
		log.Printf("[INFO] Unable to fetch the Barracuda's resource %v", err)
		return nil, err
	}

	var resourceData *WAFResouceData
	err = json.Unmarshal(data, &resourceData)

	if err != nil {
		log.Printf("[INFO] Unable to unmarshal Barracuda's resource data %v", err)
		return nil, err
	}

	return resourceData, nil
}

// deleteReq : delete APIs
func (b *BarracudaWAF) deleteReq(endpoint string) ([]byte, error) {
	req := &APIRequest{
		Method: "delete",
		URL:    endpoint,
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

// deleteReqBody : delete APIs with body
func (b *BarracudaWAF) deleteReqBody(body interface{}, endpoint string) ([]byte, error) {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "delete",
		URL:         endpoint,
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

// postReq : post APIs
func (b *BarracudaWAF) postReq(body interface{}, endpoint string) ([]byte, error) {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "post",
		URL:         endpoint,
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

// putReq : put APIs
func (b *BarracudaWAF) putReq(body interface{}, endpoint string) ([]byte, error) {
	marshalJSON, err := jsonMarshal(body)
	if err != nil {
		return nil, err
	}

	req := &APIRequest{
		Method:      "put",
		URL:         endpoint,
		Body:        strings.TrimRight(string(marshalJSON), "\n"),
		ContentType: "application/json",
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

// getReq : get APIs
func (b *BarracudaWAF) getReq(endpoint string) ([]byte, error) {
	req := &APIRequest{
		Method: "get",
		URL:    endpoint,
	}

	resp, callErr := b.APICall(req)
	return resp, callErr
}

// APICall : is used to query the Barracuda WAF web API.
func (b *BarracudaWAF) APICall(options *APIRequest) ([]byte, error) {
	var req *http.Request

	client := &http.Client{
		Transport: b.Transport,
		Timeout:   time.Second * 60,
	}

	url := fmt.Sprintf("%s/%s%s", b.Host, baseURI, options.URL)

	if options.Body != nil {
		body := bytes.NewReader([]byte(options.Body.(string)))
		req, _ = http.NewRequest(strings.ToUpper(options.Method), url, body)
	} else {
		req, _ = http.NewRequest(strings.ToUpper(options.Method), url, nil)
	}

	if b.Token != "" {
		req.Header.Set("Authorization", b.Token)
	}

	if len(options.ContentType) > 0 {
		req.Header.Set("Content-Type", options.ContentType)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		return data, b.checkError(data)
	}

	return data, nil
}

// checkError : handles any errors from API requests. It returns either the
// message of the error, if any, or nil.
func (b *BarracudaWAF) checkError(resp []byte) error {
	if len(resp) == 0 {
		return nil
	}

	var reqError RequestError

	err := json.Unmarshal(resp, &reqError)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("%s\n%s", err.Error(), string(resp[:])))
	}

	err = reqError.Error()
	if err != nil {
		return err
	}

	return nil
}

// jsonMarshal specifies an encoder with 'SetEscapeHTML' set to 'false' so that <, >, and & are not escaped.
// https://golang.org/pkg/encoding/json/#Marshal
// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
