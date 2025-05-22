// Wrapper for OVH API client just for adding rate limiting

package ovhwrap

import (
	"context"

	"github.com/ovh/go-ovh/ovh"
	"go.uber.org/ratelimit"
)

type Client struct {
	*ovh.Client
	RateLimiter ratelimit.Limiter
}

func NewClient(ovhClient *ovh.Client, rateLimiter ratelimit.Limiter) *Client {
	return &Client{
		Client:      ovhClient,
		RateLimiter: rateLimiter,
	}
}

func (c *Client) _CallAPIWithContext(ctx context.Context, method, url string, req, res interface{}, auth bool) error {
	c.RateLimiter.Take()
	return c.CallAPIWithContext(ctx, method, url, req, res, auth)
}

func (c *Client) _CallAPI(method, path string, reqBody, resType interface{}, needAuth bool) error {
	return c._CallAPIWithContext(context.Background(), method, path, reqBody, resType, needAuth)
}

func (c *Client) CallAPI(method, path string, reqBody, resType interface{}, needAuth bool) error {
	return c._CallAPI(method, path, reqBody, resType, needAuth)
}

func (c *Client) Get(url string, resType interface{}) error {
	return c._CallAPI("GET", url, nil, resType, true)
}

func (c *Client) GetUnAuth(url string, resType interface{}) error {
	return c._CallAPI("GET", url, nil, resType, false)
}

func (c *Client) Post(url string, reqBody, resType interface{}) error {
	return c._CallAPI("POST", url, reqBody, resType, true)
}

func (c *Client) PostUnAuth(url string, reqBody, resType interface{}) error {
	return c._CallAPI("POST", url, reqBody, resType, false)
}

func (c *Client) Put(url string, reqBody, resType interface{}) error {
	return c._CallAPI("PUT", url, reqBody, resType, true)
}

func (c *Client) PutUnAuth(url string, reqBody, resType interface{}) error {
	return c._CallAPI("PUT", url, reqBody, resType, false)
}

func (c *Client) Delete(url string, resType interface{}) error {
	return c._CallAPI("DELETE", url, nil, resType, true)
}

func (c *Client) DeleteUnAuth(url string, resType interface{}) error {
	return c._CallAPI("DELETE", url, nil, resType, false)
}

func (c *Client) GetWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "GET", url, nil, resType, true)
}

func (c *Client) GetUnAuthWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "GET", url, nil, resType, false)
}

func (c *Client) PostWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "POST", url, reqBody, resType, true)
}

func (c *Client) PostUnAuthWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "POST", url, reqBody, resType, false)
}

func (c *Client) PutWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "PUT", url, reqBody, resType, true)
}

func (c *Client) PutUnAuthWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "PUT", url, reqBody, resType, false)
}

func (c *Client) DeleteWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "DELETE", url, nil, resType, true)
}

func (c *Client) DeleteUnAuthWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "DELETE", url, nil, resType, false)
}
