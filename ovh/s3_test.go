package ovh

import (
	"net/http"
	"testing"
	"time"
)

// The GET Bucket Lifecycle example of the AWS documentation on Signature Version 4
// (Authenticating Requests: Using the Authorization Header).
func TestS3ClientSign(t *testing.T) {
	client := &s3Client{
		accessKeyID:     "AKIAIOSFODNN7EXAMPLE",
		secretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}

	req, err := http.NewRequest(http.MethodGet, "https://examplebucket.s3.amazonaws.com/?lifecycle", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Amz-Date", "20130524T000000Z")
	req.Header.Set("X-Amz-Content-Sha256", sha256Hex(nil))

	got := client.sign(req, "us-east-1", time.Date(2013, 5, 24, 0, 0, 0, 0, time.UTC))
	want := "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20130524/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=fea454ca298b7da1c68078a5d1bdbfbbe0d65c699e0f91ac7a200a0136783543"
	if got != want {
		t.Errorf("Authorization\n got: %s\nwant: %s", got, want)
	}
}

func TestS3AccessControlPolicyCannedACL(t *testing.T) {
	owner := s3Grantee{ID: "1542319462669586:user-5hwhM25pPT6f"}
	ownerGrant := s3Grant{Grantee: owner, Permission: "FULL_CONTROL"}

	cases := map[string]struct {
		grants []s3Grant
		want   string
	}{
		"private": {
			grants: []s3Grant{ownerGrant},
			want:   "private",
		},
		"public-read": {
			grants: []s3Grant{ownerGrant, {Grantee: s3Grantee{URI: s3GroupAllUsers}, Permission: "READ"}},
			want:   "public-read",
		},
		"public-read-write": {
			grants: []s3Grant{
				ownerGrant,
				{Grantee: s3Grantee{URI: s3GroupAllUsers}, Permission: "WRITE"},
				{Grantee: s3Grantee{URI: s3GroupAllUsers}, Permission: "READ"},
			},
			want: "public-read-write",
		},
		"authenticated-read": {
			grants: []s3Grant{ownerGrant, {Grantee: s3Grantee{URI: s3GroupAuthenticatedUsers}, Permission: "READ"}},
			want:   "authenticated-read",
		},
		"log-delivery-write": {
			grants: []s3Grant{
				ownerGrant,
				{Grantee: s3Grantee{URI: s3GroupLogDelivery}, Permission: "WRITE"},
				{Grantee: s3Grantee{URI: s3GroupLogDelivery}, Permission: "READ_ACP"},
			},
			want: "log-delivery-write",
		},
		"log-delivery-write without owner grant": {
			grants: []s3Grant{
				{Grantee: s3Grantee{URI: s3GroupLogDelivery}, Permission: "READ_ACP"},
				{Grantee: s3Grantee{URI: s3GroupLogDelivery}, Permission: "WRITE"},
			},
			want: "log-delivery-write",
		},
		"custom grant": {
			grants: []s3Grant{ownerGrant, {Grantee: s3Grantee{ID: "1542319462669586:user-other"}, Permission: "READ"}},
			want:   "",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			policy := s3AccessControlPolicy{Owner: owner, Grants: tc.grants}
			if got := policy.CannedACL(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
