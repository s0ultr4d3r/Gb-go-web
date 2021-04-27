package yadisk

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SystemFolders struct {
	Applications string `json:"applications"`
	Downloads    string `json:"downloads"`
}

type Disk struct {
	TrashSize     uint          `json:"trash_size"`
	TotalSpace    uint          `json:"total_space"`
	UsedSpace     uint          `json:"used_space"`
	SystemFolders SystemFolders `json:"system_folders"`
}

type DiskService service

func (s *DiskService) Get(ctx context.Context) (*Disk, *http.Response, error) {
	url := "disk"
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	disk := new(Disk)
	resp, err := s.client.Do(ctx, req, disk)
	if err != nil {
		return nil, resp, err
	}

	return disk, resp, nil
}

const (
	defaultBaseURL = "https://cloud-api.yandex.net/"
	apiVersion     = "1"
)

type Client struct {
	HTTPClient  *http.Client
	AccessToken string
	BaseURL     *url.URL
	Disk        *DiskService
	Resources   *ResourcesService
}

type service struct {
	client *Client
}

func NewClient(accessToken string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		HTTPClient:  http.DefaultClient,
		BaseURL:     baseURL,
		AccessToken: accessToken,
	}
	c.Disk = &DiskService{client: c}
	c.Resources = &ResourcesService{client: c}

	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	rel.Path = "v" + apiVersion + "/" + rel.Path + "/"
	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "OAuth "+c.AccessToken)

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := ctxhttp.Do(ctx, c.HTTPClient, req) // todo: test context
	if err != nil {
		return nil, err
	}

	defer func() {

		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	if err = checkResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {

			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {

			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil
			}
		}
	}

	return resp, err
}

func checkResponse(r *http.Response) error {
	if r.StatusCode >= 400 {
		apiErr := new(APIError)

		json.NewDecoder(r.Body).Decode(apiErr)

		return apiErr
	}

	return nil
}

func yaDiskSave(string) {
	var token string
	var file string
	fmt.Print("Insert your OAuth token:")
	fmt.Scan(&token)
	fmt.Print("Insert file url:")
	fmt.Scan(&file)
	fileName := path.Base(file)
	file = url.QueryEscape(file)
	cl := NewClient(token)
	req, err := NewRequest("POST", "https://cloud-api.yandex.net/v1/disk/resources/upload?path=disk%3A%2F"+fileName+file, _)
	if err != nil {
		fmt.Println(err)
	}
	Do(cl)
}
