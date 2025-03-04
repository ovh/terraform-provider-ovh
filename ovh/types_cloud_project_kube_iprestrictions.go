package ovh

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

type CloudProjectKubeIpRestrictionsCreateOrUpdateOpts struct {
	Ips []string `json:"ips"`
}

type CloudProjectKubeIpRestrictionsResponse = []string

func (opts *CloudProjectKubeIpRestrictionsCreateOrUpdateOpts) FromResource(d *schema.ResourceData) *CloudProjectKubeIpRestrictionsCreateOrUpdateOpts {
	opts.Ips, _ = helpers.StringsFromSchema(d, "ips")
	return opts
}

func (s *CloudProjectKubeIpRestrictionsCreateOrUpdateOpts) String() string {
	return fmt.Sprintf("%s", s.Ips)
}
