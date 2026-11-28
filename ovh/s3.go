package ovh

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

const s3Xmlns = "http://s3.amazonaws.com/doc/2006-03-01/"

// s3Client sends the bucket subresource requests (`?acl`, `?logging`) that OVHcloud only
// exposes through the S3 API, signed with the S3 credentials of the provider configuration.
// Signature V4 is implemented with the standard library: four bucket calls do not justify
// the AWS SDK as a dependency.
type s3Client struct {
	accessKeyID     string
	secretAccessKey string
	userAgent       string
	httpClient      *http.Client
	now             func() time.Time
}

func newS3Client(accessKeyID, secretAccessKey, userAgent string) *s3Client {
	httpClient := cleanhttp.DefaultPooledClient()
	httpClient.Transport = logging.NewTransport("OVH S3", httpClient.Transport)

	return &s3Client{
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		userAgent:       userAgent,
		httpClient:      httpClient,
		now:             time.Now,
	}
}

// s3Error is the error document S3 answers with a non 2xx status.
type s3Error struct {
	StatusCode int
	Code       string `xml:"Code"`
	Message    string `xml:"Message"`
}

func (e *s3Error) Error() string {
	return fmt.Sprintf("S3 answered %d %s: %s", e.StatusCode, e.Code, e.Message)
}

// isS3NoSuchBucket tells whether the error is S3 answering that the bucket does not exist.
func isS3NoSuchBucket(err error) bool {
	var errS3 *s3Error
	return errors.As(err, &errS3) && errS3.Code == "NoSuchBucket"
}

// do sends a signed request for a subresource of the bucket served at host, and returns the response body.
// region is the lowercase region of the bucket, as expected in the credential scope.
func (c *s3Client) do(ctx context.Context, method, host, region, subresource string, headers map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, "https://"+host+"/?"+subresource, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	now := c.now().UTC()
	req.Header.Set("X-Amz-Date", now.Format("20060102T150405Z"))
	req.Header.Set("X-Amz-Content-Sha256", sha256Hex(body))
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	req.Header.Set("Authorization", c.sign(req, region, now))
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {
		errS3 := &s3Error{StatusCode: resp.StatusCode}
		if err := xml.Unmarshal(respBody, errS3); err != nil || errS3.Code == "" {
			errS3.Code = http.StatusText(resp.StatusCode)
			errS3.Message = strings.TrimSpace(string(respBody))
		}
		return nil, errS3
	}

	return respBody, nil
}

// sign computes the Authorization header of the request (AWS Signature Version 4, s3 service).
// Every header set on the request is signed, the User-Agent is set afterwards.
func (c *s3Client) sign(req *http.Request, region string, now time.Time) string {
	names := make([]string, 0, len(req.Header)+1)
	values := map[string]string{"host": req.URL.Host}
	for name, value := range req.Header {
		name = strings.ToLower(name)
		values[name] = strings.TrimSpace(value[0])
	}
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)

	var canonicalHeaders strings.Builder
	for _, name := range names {
		canonicalHeaders.WriteString(name + ":" + values[name] + "\n")
	}
	signedHeaders := strings.Join(names, ";")

	canonicalRequest := strings.Join([]string{
		req.Method,
		req.URL.EscapedPath(),
		req.URL.RawQuery + "=",
		canonicalHeaders.String(),
		signedHeaders,
		req.Header.Get("X-Amz-Content-Sha256"),
	}, "\n")

	date := now.Format("20060102")
	scope := date + "/" + region + "/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		req.Header.Get("X-Amz-Date"),
		scope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")

	signingKey := []byte("AWS4" + c.secretAccessKey)
	for _, part := range []string{date, region, "s3", "aws4_request"} {
		signingKey = hmacSHA256(signingKey, part)
	}

	return "AWS4-HMAC-SHA256 Credential=" + c.accessKeyID + "/" + scope +
		", SignedHeaders=" + signedHeaders +
		", Signature=" + hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

// Server access logging

type s3BucketLoggingStatus struct {
	XMLName        xml.Name          `xml:"BucketLoggingStatus"`
	Xmlns          string            `xml:"xmlns,attr,omitempty"`
	LoggingEnabled *s3LoggingEnabled `xml:"LoggingEnabled"`
}

type s3LoggingEnabled struct {
	TargetBucket string `xml:"TargetBucket"`
	TargetPrefix string `xml:"TargetPrefix"`
}

// GetBucketLogging returns the logging configuration of the bucket, nil when logging is disabled.
func (c *s3Client) GetBucketLogging(ctx context.Context, host, region string) (*s3LoggingEnabled, error) {
	body, err := c.do(ctx, http.MethodGet, host, region, "logging", nil, nil)
	if err != nil {
		return nil, err
	}

	var status s3BucketLoggingStatus
	if err := xml.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("cannot decode the logging configuration of %s: %w", host, err)
	}

	return status.LoggingEnabled, nil
}

// PutBucketLogging sets the logging configuration of the bucket, nil disables logging.
func (c *s3Client) PutBucketLogging(ctx context.Context, host, region string, logging *s3LoggingEnabled) error {
	body, err := xml.Marshal(s3BucketLoggingStatus{Xmlns: s3Xmlns, LoggingEnabled: logging})
	if err != nil {
		return err
	}

	_, err = c.do(ctx, http.MethodPut, host, region, "logging", map[string]string{"Content-Type": "application/xml"}, body)
	return err
}

// Access control list

const (
	s3GroupAllUsers           = "http://acs.amazonaws.com/groups/global/AllUsers"
	s3GroupAuthenticatedUsers = "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
	s3GroupLogDelivery        = "http://acs.amazonaws.com/groups/s3/LogDelivery"
)

// s3CannedACLs lists the grants of each canned ACL, beyond the FULL_CONTROL of the bucket owner.
var s3CannedACLs = map[string][]string{
	"private":            {},
	"public-read":        {s3GroupAllUsers + ":READ"},
	"public-read-write":  {s3GroupAllUsers + ":READ", s3GroupAllUsers + ":WRITE"},
	"authenticated-read": {s3GroupAuthenticatedUsers + ":READ"},
	"log-delivery-write": {s3GroupLogDelivery + ":READ_ACP", s3GroupLogDelivery + ":WRITE"},
}

type s3AccessControlPolicy struct {
	XMLName xml.Name  `xml:"AccessControlPolicy"`
	Owner   s3Grantee `xml:"Owner"`
	Grants  []s3Grant `xml:"AccessControlList>Grant"`
}

type s3Grant struct {
	Grantee    s3Grantee `xml:"Grantee"`
	Permission string    `xml:"Permission"`
}

// s3Grantee is a canonical user (ID) or a predefined group (URI).
type s3Grantee struct {
	ID  string `xml:"ID"`
	URI string `xml:"URI"`
}

// CannedACL returns the canned ACL the grants amount to, "" when they match none of them.
// The grants of the owner are left out: every canned ACL gives the owner FULL_CONTROL.
func (p *s3AccessControlPolicy) CannedACL() string {
	grants := []string{}
	for _, grant := range p.Grants {
		if grant.Grantee.URI == "" && grant.Grantee.ID == p.Owner.ID {
			continue
		}
		grantee := grant.Grantee.URI
		if grantee == "" {
			grantee = grant.Grantee.ID
		}
		grants = append(grants, grantee+":"+grant.Permission)
	}
	sort.Strings(grants)

	for name, cannedGrants := range s3CannedACLs {
		if strings.Join(grants, ",") == strings.Join(cannedGrants, ",") {
			return name
		}
	}
	return ""
}

// GetBucketACL returns the access control list of the bucket.
func (c *s3Client) GetBucketACL(ctx context.Context, host, region string) (*s3AccessControlPolicy, error) {
	body, err := c.do(ctx, http.MethodGet, host, region, "acl", nil, nil)
	if err != nil {
		return nil, err
	}

	var policy s3AccessControlPolicy
	if err := xml.Unmarshal(body, &policy); err != nil {
		return nil, fmt.Errorf("cannot decode the ACL of %s: %w", host, err)
	}

	return &policy, nil
}

// PutBucketCannedACL sets a canned ACL on the bucket.
func (c *s3Client) PutBucketCannedACL(ctx context.Context, host, region, acl string) error {
	_, err := c.do(ctx, http.MethodPut, host, region, "acl", map[string]string{"X-Amz-Acl": acl}, nil)
	return err
}
