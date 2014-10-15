package couch

import (
	"net/http"
	"net/url"
	"io"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"fmt"
)

type Response struct {
	Ok bool
	ID string
	Rev string
	Error string
	Reason string
}

type Client struct {
	URL *url.URL
}

func NewClient(u *url.URL) *Client {
	return &Client{u}
}

func NewClientURL(urlString string) (*Client, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	return &Client{u}, nil
}

func (c *Client) AllDBs() ([]string, error) {
	res := []string{}
	_, err := c.execJSON("GET", "/_all_dbs", &res, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// TODO
//func (c *Client) AllDesignDocs() {
//}

func (c *Client) CreateDB() (resp *Response, code int, err error) {
	req, err := c.NewRequest("PUT", c.UrlString(c.DBPath(), nil), nil, nil)
	if err != nil {
		return
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	code, err = c.HandleResponse(httpResp, &resp)
	if err != nil {
		return
	}

	return 
}

func (c *Client) DeleteDB() (resp *Response, code int, err error) {
	req, err := c.NewRequest("DELETE", c.UrlString(c.DBPath(), nil), nil, nil)
	if err != nil {
		return
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	code, err = c.HandleResponse(httpResp, &resp)
	if err != nil {
		return
	}
	return 
}

func (c *Client) Save(doc interface{}) (res *Response, err error) {
	id, _, err := ParseIdRev(doc)
	if err != nil {
		return
	}

	// If no id is provided, we assume POST
	if id == "" {
		_, err = c.execJSON("POST", c.URL.Path, &res, doc, nil, nil)
	} else {
		_, err = c.execJSON("PUT", c.DocPath(id), &res, doc, nil, nil)
	}

	if err != nil {
		return
	}
	
	if res.Error != "" {
		return res, fmt.Errorf(fmt.Sprintf("%s: %s", res.Error, res.Reason))
	}

	return
}

func (c *Client) Get(id string, doc interface{}) error {
	_, err := c.execJSON("GET", c.DocPath(id), &doc, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Delete(id string, rev string) error {
	headers := http.Header{}
	headers.Add("If-Match", rev)
	res := Response{}
	_, err := c.execJSON("DELETE", c.DocPath(id), &res, nil, nil, &headers)
	if err != nil {
		return err
	}
	return nil
}

type BulkSaveRequest struct {
	Docs []interface{} `json:"docs"`
}

func (c *Client) BulkSave(docs ...interface{}) (resp *[]Response, code int, err error) {
	bulkSaveRequest := &BulkSaveRequest{Docs: docs}
	reader, err := docReader(bulkSaveRequest)
		
	req, err := c.NewRequest("POST", c.UrlString(c.DBPath() + "/_bulk_docs", nil), reader, nil)
	if err != nil {
		return
	}
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	code, err = c.HandleResponse(httpResp, &resp)
	if err != nil {
		return
	}
	return 
}

type MultiDocResponse struct {
	TotalRows uint64 `json:"total_rows"`
	Offset uint64
	Rows []Row
}

type Row struct {
	ID *string
	Key interface{}
	Value interface{}
	Doc interface{}
}

type KeysRequest struct {
	Keys []string `json:"keys"`
}

func (c *Client) View(design string, name string, options *url.Values, keys *[]string) (multiDocResponse *MultiDocResponse, err error) {
	url := c.UrlString(c.DBPath() + "/_design/" + design + "/_view/" + name, options)
	
	method := ""
	body := new(bytes.Buffer)
	if keys != nil {
		reqJson, _ := json.Marshal(KeysRequest{Keys: *keys})
		body = bytes.NewBuffer(reqJson)
		method = "POST"
	} else {
		method = "GET"
	}

	req, err := c.NewRequest(method, url, body, nil)
	if err != nil {
		return
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 
	}

	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return 
	}

	if err = json.Unmarshal(respBody, &multiDocResponse); err != nil {
		return 
	}
	return 
}

func (c *Client) Copy(src string, dest string, destRev *string) (resp *Response, code int, err error) {
	if destRev != nil {
		dest += "?rev=" + *destRev
	}

	req, err := c.NewRequest("COPY", c.UrlString(c.DocPath(src), nil), nil, nil)
	req.Header.Add("Destination", dest)
	if err != nil {
		return
	}
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	code, err = c.HandleResponse(httpResp, &resp)
	if err != nil {
		return
	}
	return 
}

type ReplicateRequest struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Cancel bool `json:"cancel,omitempty"`
	Continuous bool `json:"continuous,omitempty"`
	CreateTarget bool `json:"create_target,omitempty"`
	DocIDs []string `json:"doc_ids,omitempty"`
	Filter string `json:"filter,omitempty"`
	Proxy string `json:"proxy,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty"`
}

type ReplicateResponse struct {
	Ok bool `json:"ok"`
	LocalID bool `json:"_local_id"`
}

func (c *Client) Replicate(repReq *ReplicateRequest) (resp *ReplicateResponse, code int, err error) {
	reqReader, err := docReader(repReq)
	if err != nil {
		return
	}

	req, err := c.NewRequest("POST", c.UrlString("/_replicate", nil), reqReader, nil)
	if err != nil {
		return
	}
	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	code, err = c.HandleResponse(httpResp, &resp)
	if err != nil {
		return
	}
	return
}

func (c *Client) DBPath() string {
	return c.URL.Path
}

func (c *Client) DocPath(id string) string {
	return c.DBPath() + "/" + id
}

func (c *Client) NewRequest(method, url string, body io.Reader, headers *http.Header) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if headers != nil {
		req.Header = *headers
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return
}

func (c *Client) UrlString(path string, values *url.Values) string {
	u := *c.URL
	u.Path = path
	if values != nil {
		u.RawQuery = values.Encode()
	}
	return u.String()
}

func (c *Client) HandleResponse(resp *http.Response, result interface{}) (code int, err error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	code = resp.StatusCode
	if err = c.HandleResponseError(code, body); err != nil {
		return code, err
	}
	if err = json.Unmarshal(body, result); err != nil {
		return 0, err
	}
	return
}

func (c *Client) HandleResponseError(code int, resBytes []byte) error {
	if code < 200 || code >= 300 {
		res := Response{}
		if err := json.Unmarshal(resBytes, &res); err != nil {
			return err
		}
		return fmt.Errorf(fmt.Sprintf("Code: %d, Error: %s, Reason: %s", code, res.Error, res.Reason))
	}
	return nil
}

func (c *Client) execJSON(method string, path string, result interface{}, doc interface{}, values *url.Values, headers *http.Header) (int, error) {
	resBytes, code, err := c.execRead(method, path, doc, values, headers)
	if err != nil {
		return 0, err
	}
	if err = c.HandleResponseError(code, resBytes); err != nil {
		return code, err
	}
	if err = json.Unmarshal(resBytes, result); err != nil {
		return 0, err
	}
	return code, nil
}

func (c *Client) execRead(method string, path string, doc interface{}, values *url.Values, headers *http.Header) ([]byte, int, error) {
	r, code, err := c.exec(method, path, doc, values, headers)
	if err != nil {
		return nil, 0, err
	}

	resBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, 0, err
	}
	return resBytes, code, nil
}

func (c *Client) exec(method string, path string, doc interface{}, values *url.Values, headers *http.Header) (io.Reader, int, error) {
	reqReader, err := docReader(doc)
	if err != nil {
		return nil, 0, err
	}

	req, err := c.NewRequest(method, c.UrlString(path, values), reqReader, headers)
	if err != nil {
		return nil, 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.StatusCode, nil
}

func docReader(doc interface{}) (io.Reader, error) {
	if doc == nil {
		return nil, nil
	}

	docJson, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}
	r := bytes.NewBuffer(docJson)
	return r, nil
}

type IdRev struct {
	ID string `json:"_id"`
	Rev string `json:"_rev"`
}

func Remarshal(doc interface{}, newDoc interface{}) (err error) {
	docJson, err := json.Marshal(doc)
	if err != nil {
		return 
	}

	err = json.Unmarshal(docJson, newDoc)
	if err != nil {
		return 
	}
	return
}

func ParseIdRev(doc interface{}) (string, string, error) {
	docJson, err := json.Marshal(doc)
	if err != nil {
		return "", "", err
	}
	
	idRev := &IdRev{}
	if err = json.Unmarshal(docJson, idRev); err != nil {
		return "", "", err
	}

	return idRev.ID, idRev.Rev, nil
}

