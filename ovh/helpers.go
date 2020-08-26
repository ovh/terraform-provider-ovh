package ovh

import (
	"bytes"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func validateBootType(value string) error {
	return validateStringEnum(value, []string{
		"harddisk",
		"internal",
		"ipxeCustomerScript",
		"network",
		"rescue",
	})
}

func validateLanguageCode(value string) error {
	// accepted language code for dedicated servers
	return validateStringEnum(value, []string{
		"ar",
		"bg",
		"cs",
		"da",
		"de",
		"el",
		"en",
		"es",
		"et",
		"fi",
		"fr",
		"he",
		"hr",
		"hu",
		"it",
		"ja",
		"ko",
		"lt",
		"lv",
		"nb",
		"nl",
		"no",
		"pl",
		"pt",
		"ro",
		"ru",
		"sk",
		"sl",
		"sr",
		"sv",
		"th",
		"tr",
		"tu",
		"uk",
		"zh-Hans-CN",
		"zh-Hans-HK",
	})
}

func validateRAIDMode(value string) error {
	// accepted raid modes for installation templates hardware specs
	return validateStringEnum(value, []string{
		"raid0",
		"raid1",
		"raid10",
		"raid5",
		"raid50",
		"raid6",
		"raid60",
	})
}

func validatePartitionType(value string) error {
	// accepted partition types for installation templates
	return validateStringEnum(value, []string{
		"lv",
		"primary",
		"logical",
	})
}

func validatePartitionRAIDMode(value string) error {
	// accepted raid modes for installation templates partitions specs
	return validateStringEnum(value, []string{
		"raid0",
		"raid1",
		"raid10",
		"raid5",
		"raid6",
	})
}

func validateFilesystem(value string) error {
	// accepted filesystem types for installation templates partitions specs
	return validateStringEnum(value, []string{
		"btrfs",
		"ext3",
		"ext4",
		"ntfs",
		"reiserfs",
		"swap",
		"ufs",
		"xfs",
		"zfs",
	})
}

func validateDedicatedCephCrushTunables(value string) error {
	return validateStringEnum(value, []string{
		"OPTIMAL",
		"DEFAULT",
		"LEGACY",
		"BOBTAIL",
		"ARGONAUT",
		"FIREFLY",
		"HAMMER",
		"JEWEL",
	})
}

func validateDedicatedCephStatus(value string) error {
	return validateStringEnum(value, []string{
		"CREATING",
		"INSTALLED",
		"DELETING",
		"DELETED",
		"TASK_IN_PROGRESS",
	})
}

func validateDedicatedCephACLFamily(value string) error {
	return validateStringEnum(value, []string{
		"IPv4",
		"IPv6",
	})
}

func getNilBoolPointerFromData(data interface{}, id string) *bool {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		if val, ok := resourceData.GetOk(id); ok {
			return getNilBoolPointer(val)
		}
	} else if mapData, tok := data.(map[string]interface{}); tok {
		if val, ok := mapData[id]; ok {
			return getNilBoolPointer(val)
		}
	}

	return nil
}

func getNilStringPointerFromData(data interface{}, id string) *string {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		if val, ok := resourceData.GetOk(id); ok {
			return getNilStringPointer(val)
		}
	} else if mapData, tok := data.(map[string]interface{}); tok {
		if val, ok := mapData[id]; ok {
			return getNilStringPointer(val)
		}
	}

	return nil
}

func getNilIntPointerFromData(data interface{}, id string) *int {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		if val, ok := resourceData.GetOk(id); ok {
			return getNilIntPointer(val)
		}
	} else if mapData, tok := data.(map[string]interface{}); tok {
		if val, ok := mapData[id]; ok {
			return getNilIntPointer(val)
		}
	}

	return nil
}

func getNilInt64PointerFromData(data interface{}, id string) *int64 {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		if val, ok := resourceData.GetOk(id); ok {
			return getNilInt64Pointer(val)
		}
	} else if mapData, tok := data.(map[string]interface{}); tok {
		if val, ok := mapData[id]; ok {
			return getNilInt64Pointer(val)
		}
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

func getNilInt64Pointer(val interface{}) *int64 {
	if val == nil {
		return nil
	}
	value := int64(val.(int))
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
	xs := []string{}
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
