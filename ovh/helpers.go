package ovh

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func validateIpBlock(value string) error {
	if _, _, err := net.ParseCIDR(value); err != nil {
		return fmt.Errorf("Value %s is not a valid IP block", value)
	}
	return nil
}

func validateIp(value string) error {
	if ip := net.ParseIP(value); ip != nil {
		return nil
	}
	return fmt.Errorf("Value %s is not a valid IP", value)
}

func validateIpV6(value string) error {
	if ip := net.ParseIP(value); ip != nil && ip.To4() == nil && ip.To16() != nil {
		return nil
	}
	return fmt.Errorf("Value %s is not a valid IPv6", value)
}

func validateIpV4(value string) error {
	if ip := net.ParseIP(value); ip != nil && ip.To4() != nil {
		return nil
	}
	return fmt.Errorf("Value %s is not a valid IPv4", value)
}

func validateStringEnum(value string, enum []string) error {
	missing := true
	for _, v := range enum {
		if value == v {
			missing = false
		}
	}
	if missing {
		return fmt.Errorf("Value %s is not among valid values (%s)", value, enum)
	}
	return nil
}

func getNilBoolPointer(val interface{}) *bool {
	if val == nil {
		return nil
	}
	value := val.(bool)
	return &value
}

func getNilStringPointer(val interface{}) *string {
	if val == nil {
		return nil
	}
	value := val.(string)
	if len(value) == 0 {
		return nil
	}
	return &value
}

func getNilIntPointer(val interface{}) *int {
	if val == nil {
		return nil
	}
	value := val.(int)
	return &value
}

func conditionalAttributeInt(buff *bytes.Buffer, name string, val *int) {
	if val != nil {
		buff.WriteString(fmt.Sprintf("  %s = %d\n", name, *val))
	}
}

func conditionalAttributeString(buff *bytes.Buffer, name string, val *string) {
	if val != nil {
		buff.WriteString(fmt.Sprintf("  %s = \"%s\"\n", name, *val))
	}
}

func conditionalAttributeBool(buff *bytes.Buffer, name string, val *bool) {
	if val != nil {
		buff.WriteString(fmt.Sprintf("  %s = %v\n", name, *val))
	}
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, endpoint string) error {
	if err.(*ovh.APIError).Code == 404 {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
}

func stringsFromSchema(d *schema.ResourceData, id string) []string {
	var xs []string
	if v := d.Get(id); v != nil {
		rs := v.(*schema.Set).List()
		if len(rs) > 0 {
			for _, v := range v.(*schema.Set).List() {
				xs = append(xs, v.(string))
			}
		}
	}
	return xs
}

func normalizeIPSubnet(ip string) (string, error) {
	if !strings.Contains(ip, "/") {
		ip += "/32"
	}
	return ip, nil
}
