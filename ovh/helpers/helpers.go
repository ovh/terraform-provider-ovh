package helpers

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func ValidateIpBlock(value string) error {
	if _, _, err := net.ParseCIDR(value); err != nil {
		return fmt.Errorf("Value %s is not a valid IP block", value)
	}
	return nil
}

func ValidateIp(value string) error {
	if ip := net.ParseIP(value); ip != nil {
		return nil
	}
	return fmt.Errorf("Value %s is not a valid IP", value)
}

func ValidateIpV6(value string) error {
	if ip := net.ParseIP(value); ip != nil && ip.To4() == nil && ip.To16() != nil {
		return nil
	}
	return fmt.Errorf("Value %s is not a valid IPv6", value)
}

func ValidateIpV4(value string) error {
	if ip := net.ParseIP(value); ip != nil && ip.To4() != nil {
		return nil
	}
	return fmt.Errorf("Value %s is not a valid IPv4", value)
}

func ValidateStringEnum(value string, enum []string) error {
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

func ValidateBootType(value string) error {
	return ValidateStringEnum(value, []string{
		"harddisk",
		"internal",
		"ipxeCustomerScript",
		"network",
		"rescue",
	})
}

func ValidateLanguageCode(value string) error {
	// accepted language code for dedicated servers
	return ValidateStringEnum(value, []string{
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

func ValidateRAIDMode(value string) error {
	// accepted raid modes for installation templates hardware specs
	return ValidateStringEnum(value, []string{
		"raid0",
		"raid1",
		"raid10",
		"raid5",
		"raid50",
		"raid6",
		"raid60",
	})
}

func ValidatePartitionType(value string) error {
	// accepted partition types for installation templates
	return ValidateStringEnum(value, []string{
		"lv",
		"primary",
		"logical",
	})
}

func ValidatePartitionRAIDMode(value string) error {
	// accepted raid modes for installation templates partitions specs
	return ValidateStringEnum(value, []string{
		"raid0",
		"raid1",
		"raid10",
		"raid5",
		"raid6",
	})
}

func ValidateFilesystem(value string) error {
	// accepted filesystem types for installation templates partitions specs
	return ValidateStringEnum(value, []string{
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

func ValidateDedicatedCephCrushTunables(value string) error {
	return ValidateStringEnum(value, []string{
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

func ValidateDedicatedCephStatus(value string) error {
	return ValidateStringEnum(value, []string{
		"CREATING",
		"INSTALLED",
		"DELETING",
		"DELETED",
		"TASK_IN_PROGRESS",
	})
}

func ValidateDedicatedCephACLFamily(value string) error {
	return ValidateStringEnum(value, []string{
		"IPv4",
		"IPv6",
	})
}

func GetNilBoolPointerFromData(data interface{}, id string) *bool {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		return GetNilBoolPointer(resourceData.Get(id).(bool))
	} else if mapData, tok := data.(map[string]interface{}); tok {
		return GetNilBoolPointer(mapData[id].(bool))
	}

	return nil
}

func GetNilStringPointerFromData(data interface{}, id string) *string {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		if val, ok := resourceData.GetOk(id); ok {
			return GetNilStringPointer(val)
		}
	} else if mapData, tok := data.(map[string]interface{}); tok {
		if val, ok := mapData[id]; ok {
			return GetNilStringPointer(val)
		}
	}

	return nil
}

func GetNilIntPointerFromData(data interface{}, id string) *int {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		if val, ok := resourceData.GetOk(id); ok {
			return GetNilIntPointer(val)
		}
	} else if mapData, tok := data.(map[string]interface{}); tok {
		if val, ok := mapData[id]; ok {
			return GetNilIntPointer(val)
		}
	}

	return nil
}

func GetNilInt64PointerFromData(data interface{}, id string) *int64 {
	if resourceData, tok := data.(*schema.ResourceData); tok {
		if val, ok := resourceData.GetOk(id); ok {
			return GetNilInt64Pointer(val)
		}
	} else if mapData, tok := data.(map[string]interface{}); tok {
		if val, ok := mapData[id]; ok {
			return GetNilInt64Pointer(val)
		}
	}

	return nil
}

func GetNilBoolPointer(value bool) *bool {
	return &value
}

func GetNilStringPointer(val interface{}) *string {
	if val == nil {
		return nil
	}
	value := val.(string)
	if len(value) == 0 {
		return nil
	}
	return &value
}

func GetNilIntPointer(val interface{}) *int {
	if val == nil {
		return nil
	}
	value := val.(int)
	return &value
}

func GetNilInt64Pointer(val interface{}) *int64 {
	if val == nil {
		return nil
	}
	value := int64(val.(int))
	return &value
}

func ConditionalAttributeInt(buff *bytes.Buffer, name string, val *int) {
	if val != nil {
		buff.WriteString(fmt.Sprintf("  %s = %d\n", name, *val))
	}
}

func ConditionalAttributeString(buff *bytes.Buffer, name string, val *string) {
	if val != nil {
		buff.WriteString(fmt.Sprintf("  %s = \"%s\"\n", name, *val))
	}
}

func ConditionalAttributeBool(buff *bytes.Buffer, name string, val *bool) {
	if val != nil {
		buff.WriteString(fmt.Sprintf("  %s = %v\n", name, *val))
	}
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, endpoint string) error {
	if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("calling %s:\n\t %s", endpoint, err.Error())
}

func StringsFromSchema(d *schema.ResourceData, id string) ([]string, error) {
	xs := []string{}
	if v := d.Get(id); v != nil {
		switch vv := v.(type) {
		case *schema.Set:
			rs := vv.List()
			if len(rs) > 0 {
				for _, vvv := range vv.List() {
					xs = append(xs, vvv.(string))
				}
			}
		case []interface{}:
			if len(vv) > 0 {
				for _, vvv := range vv {
					xs = append(xs, vvv.(string))
				}
			}
		default:
			return nil, fmt.Errorf("Attribute %v is not a list or set", id)
		}
	}
	return xs, nil
}

// WaitAvailable wait for a ressource to become available in the API (aka non 404)
func WaitAvailable(client *ovh.Client, endpoint string, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		if err := client.Get(endpoint, nil); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
}

func ValidateSubsidiary(v string) error {
	return ValidateStringEnum(strings.ToLower(v), []string{
		"cz",
		"de",
		"es",
		"eu",
		"fi",
		"fr",
		"gb",
		"ie",
		"it",
		"lt",
		"ma",
		"nl",
		"pl",
		"pt",
		"sn",
		"tn",
	})
}
