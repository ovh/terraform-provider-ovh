// Wrapper for OVH API client just for adding rate limiting

package ovh

import (
	"context"

	"github.com/ovh/go-ovh/ovh"
	"go.uber.org/ratelimit"
)

type OVHClient struct {
	RawOVHClient *ovh.Client
	RateLimiter  ratelimit.Limiter
}

func (c *OVHClient) _CallAPIWithContext(ctx context.Context, method, url string, req, res interface{}, auth bool) error {
	c.RateLimiter.Take()
	return c.RawOVHClient.CallAPIWithContext(ctx, method, url, req, res, auth)
}

func (c *OVHClient) _CallAPI(method, path string, reqBody, resType interface{}, needAuth bool) error {
	return c._CallAPIWithContext(context.Background(), method, path, reqBody, resType, needAuth)
}

func (c *OVHClient) CallAPI(method, path string, reqBody, resType interface{}, needAuth bool) error {
	return c._CallAPI(method, path, reqBody, resType, needAuth)
}

func (c *OVHClient) Get(url string, resType interface{}) error {
	return c._CallAPI("GET", url, nil, resType, true)
}

func (c *OVHClient) GetUnAuth(url string, resType interface{}) error {
	return c._CallAPI("GET", url, nil, resType, false)
}

func (c *OVHClient) Post(url string, reqBody, resType interface{}) error {
	return c._CallAPI("POST", url, reqBody, resType, true)
}

func (c *OVHClient) PostUnAuth(url string, reqBody, resType interface{}) error {
	return c._CallAPI("POST", url, reqBody, resType, false)
}

func (c *OVHClient) Put(url string, reqBody, resType interface{}) error {
	return c._CallAPI("PUT", url, reqBody, resType, true)
}

func (c *OVHClient) PutUnAuth(url string, reqBody, resType interface{}) error {
	return c._CallAPI("PUT", url, reqBody, resType, false)
}

func (c *OVHClient) Delete(url string, resType interface{}) error {
	return c._CallAPI("DELETE", url, nil, resType, true)
}

func (c *OVHClient) DeleteUnAuth(url string, resType interface{}) error {
	return c._CallAPI("DELETE", url, nil, resType, false)
}

func (c *OVHClient) GetWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "GET", url, nil, resType, true)
}

func (c *OVHClient) GetUnAuthWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "GET", url, nil, resType, false)
}

func (c *OVHClient) PostWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "POST", url, reqBody, resType, true)
}

func (c *OVHClient) PostUnAuthWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "POST", url, reqBody, resType, false)
}

func (c *OVHClient) PutWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "PUT", url, reqBody, resType, true)
}

func (c *OVHClient) PutUnAuthWithContext(ctx context.Context, url string, reqBody, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "PUT", url, reqBody, resType, false)
}

func (c *OVHClient) DeleteWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "DELETE", url, nil, resType, true)
}

func (c *OVHClient) DeleteUnAuthWithContext(ctx context.Context, url string, resType interface{}) error {
	return c._CallAPIWithContext(ctx, "DELETE", url, nil, resType, false)
}

func (c *OVHClient) Endpoint() string {
	return c.RawOVHClient.Endpoint()
}
