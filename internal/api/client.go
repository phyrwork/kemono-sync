package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *resty.Client
}

type CreatorPostsParams struct {
	Domain  string
	Service string
	User    string
}

func (c *Client) GetCreatorPosts(params CreatorPostsParams) (*[]Post, error) {
	url := fmt.Sprintf("https://%s/api/v1/%s/user/%s", params.Domain, params.Service, params.User)
	resp, err := c.client.R().SetResult(&[]Post{}).Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting posts: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("error getting posts: %v", resp.Status())
	}
	return resp.Result().(*[]Post), nil
}

type PostParams struct {
	Domain  string
	Service string
	User    string
	ID      string
}

func (c *Client) GetPost(params PostParams) (*Post, error) {
	url := fmt.Sprintf("https://%s/api/v1/%s/user/%s/post/%s", params.Domain, params.Service, params.User, params.ID)
	resp, err := c.client.R().SetResult(&Post{}).Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting post: %w", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("error getting post: %v", resp.Status())
	}
	return resp.Result().(*Post), nil
}

type DataParams struct {
	Domain string
	Path   string
}

func (c *Client) GetData(params DataParams) (io.ReadCloser, error) {
	resp, err := c.client.R().SetDoNotParseResponse(true).Get(fmt.Sprintf("https://%s/data%s", params.Domain, params.Path))
	if err != nil {
		return nil, fmt.Errorf("error getting file: %v", err)
	}
	return resp.RawBody(), nil
}

func New() *Client {
	t := &http.Transport{
		// The servers can be very slow to connect a lot of the time.
		TLSHandshakeTimeout: 60 * time.Second,
	}
	c := &http.Client{Transport: t}
	return &Client{
		client: resty.NewWithClient(c).SetRetryCount(5).SetRetryWaitTime(250 * time.Millisecond),
	}
}
