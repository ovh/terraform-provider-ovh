package ovh

import (
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// RegionAttributesHash creates an hash for the region attributes.
func RegionAttributesHash(v interface{}) int {
	attributes, ok := v.(map[string]interface{})
	if !ok {
		err := errors.New("bad casting of region attributes")
		log.Printf("[ERROR] %s: %v ", err, v)

		panic(err)
	}

	builder := strings.Builder{}

	for _, key := range []string{"status", "region", "openstackid"} {
		value, inMap := attributes[key]
		if inMap {
			stringValue, ok := value.(string)
			if ok {
				builder.WriteString(stringValue)
			} else {
				err := errors.New("bad casting of value in region attributes")
				log.Printf("[ERROR] %s on key %s with current value: %v ", err, key, value)

				panic(err)
			}
		}
	}

	return schema.HashString(builder.String())
}
